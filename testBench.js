// bench.js — 压测 Broker 代理的并发能力
/**
   $env:NODE_PATH = (npm root -g)
   node .\testBench.js --conns=1000 --rate=10 --duration=30 --mode=mixed --worldRatio=0.1
**/

const WebSocket = require('ws');
const http = require('http');
const { performance } = require('perf_hooks');
const protobuf = require('protobufjs');
const path = require('path');

// ── Protobuf 动态加载 ──
const pbRoot = protobuf.loadSync(path.join(__dirname, 'proto', 'ws.proto'));
const WSReq  = pbRoot.lookupType('ws.WSReq');
const WSResp = pbRoot.lookupType('ws.WSResp');

// ── name → oneof 字段名映射表 ──
// protobufjs 动态加载会把 proto snake_case 自动转成 camelCase
const nameToOneofField = {
    'account.login':     'accountLogin',
    'account.register':  'accountRegister',
    'chat.world.send':   'chatWorldSend',
    'chat.private.send': 'chatPrivateSend',
    'chat.group.send':   'chatGroupSend',
    'chat.group.join':   'groupMemberReq',
    'chat.group.leave':  'groupMemberReq',
    'chat.group.members':'groupMemberReq',
    'room.create':       'roomCreate',
    'room.join':         'roomJoin',
    'room.leave':        'roomLeave',
    'room.kick':         'roomKick',
    'room.delete':       'roomDelete',
    'room.send':         'roomSend',
    'inventory.add':     'inventoryAdd',
    'inventory.consume': 'inventoryConsume',
    'inventory.drop':    'inventoryDrop',
    'inventory.get':     'inventoryGet',
};

const respOneofFieldToName = {
    'accountLoginResp':    'account.login',
    'accountRegisterResp': 'account.register',
    'chatWorldNew':        'chat.world.new',
    'chatPrivateNew':      'chat.private.new',
    'chatGroupNew':        'chat.group.new',
    'groupMemberResp':     'chat.group.members',
    'chatOnlineResp':      'chat.online.get',
    'roomCreateResp':      'room.create',
    'roomJoinResp':        'room.join',
    'roomChatNew':         'room.chat.new',
    'roomNotice':          'room.notice',
    'systemBroadcast':     'system.broadcast',
    'systemPush':          'system.push',
    'inventoryAddResp':     'inventory.add',
    'inventoryConsumeResp': 'inventory.consume',
    'inventoryDropResp':    'inventory.drop',
    'inventoryGetResp':     'inventory.get',
};

// 发消息：{seq, name, msg} → WSReq encode → Buffer
// ✅ 阶段 2：优先 oneof typed message；兜底 msg_json
function encodeMsg(seq, name, msg = {}) {
    const payload = { seq: seq || 0, name };

    const oneofField = nameToOneofField[name];
    if (oneofField) {
        payload[oneofField] = msg; // oneof typed message
    } else {
        payload.msgJson = JSON.stringify(msg); // 兜底
    }

    return WSReq.encode(WSReq.create(payload)).finish();
}

// 收消息：Buffer → WSResp / BatchResp → 返回数组（统一接口）
// 阶段 3：支持服务端合并多条成 BatchResp
function decodeMsg(buf) {
    // 先尝试 BatchResp
    try {
        const BR = pbRoot.lookupType('ws.BatchResp');
        const p = BR.decode(buf);
        if (p.msgs && p.msgs.length > 0) {
            return p.msgs.map(m => unwrapWSResp(m));
        }
    } catch {}
    // 单条 WSResp
    return [unwrapWSResp(WSResp.decode(buf))];
}

// 从 WSResp protobufjs object 里提取业务 name + msg
function unwrapWSResp(p) {
    let msg = null;
    let oneofName = null;
    // 找 oneof 字段（protobufjs oneof actual 是字段名但驼峰）
    const oneofs = WSResp.oneofs;
    for (const [oofName, oof] of Object.entries(oneofs)) {
        const fieldName = p[oof]; // oneof actual field name
        if (fieldName) {
            // 遍历 oneof 的 members 找哪个被设置了
            for (const member of oof.oneof) {
                if (p[member] != null) {
                    msg = p[member];
                    oneofName = respOneofFieldToName[member];
                    break;
                }
            }
            break;
        }
    }

    // 兜底 msg_json
    if (msg == null && p.msgJson) {
        try { msg = JSON.parse(p.msgJson); } catch {}
    }

    return {
        seq: Number(p.seq),
        name: oneofName || p.name,
        code: p.code,
        msg: msg || {}
    };
}

const args = {
    conns: 1000, rate: 10, duration: 30,
    url: 'ws://127.0.0.1:9001/ws', dispatch: false,
    mode: 'mixed',
    // mixed 模式各协议比例（加起来 1.0，剩下的给 inventory）
    worldRatio: 0.2, privateRatio: 0.2, roomRatio: 0.2, groupRatio: 0.1,
    // room 模式参数
    roomSize: 50,
    // group 模式参数
    groupId: 'bench_group_01',
    // inventory 模式参数
    itemIds: 'gold,sword,potion'
};
process.argv.slice(2).forEach(a => {
    const [k, v] = a.split('=');
    const key = k.replace(/^--/, '');
    if (args[key] !== undefined) args[key] = isNaN(v) ? v : +v;
});

const stats = {
    connected: 0,
    sent: 0, sentWorld: 0, sentPrivate: 0,
    sentRoom: 0, sentGroup: 0, sentInventory: 0,
    recv: 0,
    errors: 0,
    loginOk: 0, loginFail: 0,
    connectTime: 0,
    startAt: 0,
    // 分协议精确计数
    sentByProtocol: {},
    // —— 新增：错误分类统计 ——
    errorsByProtocol: {},   // 协议名 → 错误数
    errorsByCode: {},       // 错误码 → 错误数
    errorSamples: [],       // 前 20 条错误样本（避免刷屏）
    // —— 新增：连接生命周期 ——
    disconnectCount: 0,
    closeReasons: {},       // close code → 次数
    loginDuration: 0,
};

// sentByProtocol 预初始化所有已知协议
[   'account.login', 'account.register',
    'chat.world.send', 'chat.private.send',
    'chat.group.send', 'chat.group.join', 'chat.group.leave', 'chat.group.members',
    'chat.online.get',
    'room.create', 'room.join', 'room.leave', 'room.kick', 'room.delete', 'room.send',
    'inventory.add', 'inventory.consume', 'inventory.drop', 'inventory.get'
].forEach(n => stats.sentByProtocol[n] = 0);

// uidMap: idx → 后端返回的真实 uid（如 uid_a1b2c3）
// private.send 必须用真实 uid，Redis GetUserNode 才能查到目标方是否在线
const uidMap = new Map();

let globalSeq = 0;
function nextSeq() { return ++globalSeq; }

function printStats() {
    const now = performance.now();
    const elapsed = (now - stats.startAt) / 1000;
    const sendQPS = (stats.sent / elapsed).toFixed(0);
    const recvQPS = (stats.recv / elapsed).toFixed(0);
    const recvRatio = stats.sent > 0 ? ((stats.recv / stats.sent) * 100).toFixed(1) : 0;
    console.log(`\r[${elapsed.toFixed(0)}s] conn=${stats.connected} login=${stats.loginOk}✅/${stats.loginFail}❌ sent=${stats.sent} (${sendQPS}/s) recv=${stats.recv} (${recvQPS}/s) ratio=${recvRatio}% err=${stats.errors}`, { clear: true });
}

async function resolveBrokerURL() {
    if (!args.dispatch) return args.url;
    return new Promise((resolve, reject) => {
        const url = new URL(args.url);
        const gwHost = `http://${url.hostname}:9000`;
        http.get(`${gwHost}/gateway/dispatch?uid=bench_user`, res => {
            let d = '';
            res.on('data', c => d += c);
            res.on('end', () => { try { resolve(JSON.parse(d).broker_url || args.url); } catch { resolve(args.url); } });
        }).on('error', reject);
    });
}

async function run() {
    const brokerURL = await resolveBrokerURL();
    console.log(`\n🚀 压测开始`);
    console.log(`   目标: ${brokerURL}`);
    console.log(`   连接数: ${args.conns}  速率: ${args.rate}  模式: ${args.mode}${args.mode==='mixed'?`(world${args.worldRatio*100}%)`:''}`);
    console.log(`   持续: ${args.duration}s\n`);

    stats.startAt = performance.now();
    const connectStart = performance.now();
    const sendInterval = 1000 / args.rate;

    const batchSize = 500;
    for (let i = 0; i < args.conns; i += batchSize) {
        const batch = Math.min(batchSize, args.conns - i);
        await Promise.all(Array.from({ length: batch }, (_, j) => connect(brokerURL, i + j, sendInterval)));
    }

    stats.connectTime = ((performance.now() - connectStart) / 1000).toFixed(2);
    console.log(`\n✅ ${args.conns} 个连接建立完成，耗时 ${stats.connectTime}s`);
    console.log(`📊 实时统计（Ctrl+C 提前结束）:`);

    setInterval(printStats, 2000);

    setTimeout(() => { summary(); process.exit(0); }, args.duration * 1000);
}

function connect(url, idx, interval) {
    return new Promise(resolve => {
        const ws = new WebSocket(url, { perMessageDeflate: false, handshakeTimeout: 10000 });
        let ready = false, sendTimer = null, loginDone = false, loginRetries = 0;
        const MAX_RETRIES = 2;

        // ── 模式路由 ──
        let sendMode = args.mode;
        if (args.mode === 'mixed') {
            // 按比例选一种：world / private / room / group / inventory
            const r = Math.random();
            let acc = 0;
            if ((acc += (args.worldRatio || 0.2)) > r) sendMode = 'world';
            else if ((acc += (args.privateRatio || 0.2)) > r) sendMode = 'private';
            else if ((acc += (args.roomRatio || 0.2)) > r) sendMode = 'room';
            else if ((acc += (args.groupRatio || 0.1)) > r) sendMode = 'group';
            else sendMode = 'inventory';
        }

        // ── 场景前置变量 ──
        // 注意：room.create 传的是 name（房间名），后端返回真实 room_id（rm_xxxxxxxx）
        // room.join / room.send 必须用真实 room_id，不能用 name
        const myRoomName = `bench_room_${Math.floor(idx / (args.roomSize || 50))}`;
        let myRealRoomId = '';  // room.create 响应回来后赋值
        const myGroupId = args.groupId || 'bench_group_01';
        const itemIds = (args.itemIds || 'gold,sword,potion').split(',');
        let roomCreated = false, roomJoined = false;
        let groupJoined = false, groupJoinFailed = false;
        let invRound = 0;

        const username = `bench_${idx}`;
        const password = '123456';

        ws.on('open', () => {
            stats.connected++; ready = true;
            ws.send(encodeMsg(nextSeq(), 'account.login', { username, password }));
            resolve();
        });

        ws.on('message', (raw) => {
            // Protobuf decode（支持 BatchResp / 单条 WSResp，返回数组）
            const msgs = decodeMsg(raw);
            stats.recv += msgs.length;

            for (const data of msgs) {
                if (stats.recv <= 10) {
                    console.log(`  [DEBUG recv #${stats.recv}] ${JSON.stringify(data).substring(0, 300)}`);
                }

                // —— 所有消息都做错误码统计（不管 loginDone 与否）——
                if (data.code !== undefined && data.code !== 0 && data.code !== null) {
                    stats.errors++;
                    const en = data.name || 'unknown';
                    stats.errorsByProtocol[en] = (stats.errorsByProtocol[en] || 0) + 1;
                    const ec = String(data.code);
                    stats.errorsByCode[ec] = (stats.errorsByCode[ec] || 0) + 1;
                    if (stats.errorSamples.length < 20) {
                        stats.errorSamples.push({
                            idx, name: en, code: data.code,
                            msg: JSON.stringify(data.msg).substring(0, 120)
                        });
                    }
                }

                if (loginDone) {
                    // —— 前置协议响应处理（room/group 模式）——

                    // room 模式 阶段 1：create 响应 → 存真实 roomId
                    if (sendMode === 'room' && !myRealRoomId && data.name === 'room.create') {
                        if (data.code === 0 && data.msg) {
                            // proto 里字段是 room_id，protobufjs decode 后驼峰成 roomId
                            myRealRoomId = data.msg.roomId || data.msg.room_id || '';
                        }
                        // 不管成功失败，create 只发一次，等下一轮 sendTimer 会检测 myRealRoomId
                        continue;
                    }

                    // room 模式 阶段 2：join 响应 → 设 roomJoined=true
                    if (sendMode === 'room' && !roomJoined && data.name === 'room.join') {
                        if (data.code === 0) {
                            roomJoined = true;
                            ws.send(encodeMsg(nextSeq(), 'chat.online.get', {}));
                            stats.sent++; stats.sentByProtocol['chat.online.get']++;
                        }
                        // join 失败也不重试（房间满了/不存在都是服务端问题）
                        continue;
                    }

                    // group 模式：join 响应 → 设 groupJoined=true
                    if (sendMode === 'group' && !groupJoined && data.name === 'chat.group.join') {
                        if (data.code === 0) {
                            groupJoined = true;
                        } else {
                            groupJoinFailed = true;
                        }
                        continue;
                    }

                    continue;
                }

                try {
                    if (data.name === 'account.login') {
                        if (data.code === 0 && data.msg && data.msg.token) {
                            loginDone = true;
                            stats.loginOk++;
                            // ✅ 存真实 uid（后端返回格式: uid_xxx，不是 username bench_X）
                            // protobufjs 动态解码字段名是 camelCase: uId
                            const realUid = data.msg.uId || data.msg.uid || data.msg.UId || '';
                            if (realUid) uidMap.set(idx, realUid);
                            const tgt = (idx + 1) % args.conns;
                            const sendOne = () => {
                                if (ws.readyState !== WebSocket.OPEN) return;
                                invRound++;

                                switch (sendMode) {
                                    // ──── world 世界频道（现有）────
                                    case 'world': {
                                        ws.send(encodeMsg(nextSeq(), 'chat.world.send', { content: `bench_${idx}_${Date.now()%100000}` }));
                                        stats.sent++; stats.sentWorld++;
                                        stats.sentByProtocol['chat.world.send']++;
                                        break;
                                    }
                                    // ──── private 私聊（现有）────
                                    case 'private': {
                                        const tgtUid = uidMap.get(tgt);
                                        if (!tgtUid) return;
                                        ws.send(encodeMsg(nextSeq(), 'chat.private.send', { toUid: tgtUid, content: `hi_${idx}` }));
                                        stats.sent++; stats.sentPrivate++;
                                        stats.sentByProtocol['chat.private.send']++;
                                        break;
                                    }
                                    // ──── room 房间（新增）────
                                    case 'room': {
                                        if (!roomCreated) {
                                            // 阶段 1：发 create，等响应拿真实 roomId
                                            ws.send(encodeMsg(nextSeq(), 'room.create', { name: myRoomName, maxSize: args.roomSize || 50 }));
                                            stats.sent++; stats.sentRoom++;
                                            stats.sentByProtocol['room.create']++;
                                            roomCreated = true;  // 只发一次 create
                                        } else if (!roomJoined && myRealRoomId) {
                                            // 阶段 2：create 响应已回来，用真实 roomId 发 join
                                            ws.send(encodeMsg(nextSeq(), 'room.join', { roomId: myRealRoomId }));
                                            stats.sent++; stats.sentRoom++;
                                            stats.sentByProtocol['room.join']++;
                                        } else if (roomJoined) {
                                            // 阶段 3：已经在房间里了，正常发消息
                                            ws.send(encodeMsg(nextSeq(), 'room.send', { roomId: myRealRoomId, content: `room_${idx}_${Date.now()%100000}` }));
                                            stats.sent++; stats.sentRoom++;
                                            stats.sentByProtocol['room.send']++;
                                        }
                                        break;
                                    }
                                    // ──── group 群组（新增）────
                                    case 'group': {
                                        if (!groupJoined) {
                                            ws.send(encodeMsg(nextSeq(), 'chat.group.join', { groupId: myGroupId }));
                                            stats.sent++; stats.sentGroup++;
                                            stats.sentByProtocol['chat.group.join']++;
                                            groupJoined = true;
                                        } else {
                                            ws.send(encodeMsg(nextSeq(), 'chat.group.send', { groupId: myGroupId, content: `group_${idx}_${Date.now()%100000}` }));
                                            stats.sent++; stats.sentGroup++;
                                            stats.sentByProtocol['chat.group.send']++;
                                            // 每 10 轮查一次群成员
                                            if (invRound % 10 === 0) {
                                                ws.send(encodeMsg(nextSeq(), 'chat.group.members', { groupId: myGroupId }));
                                                stats.sent++; stats.sentGroup++;
                                                stats.sentByProtocol['chat.group.members']++;
                                            }
                                        }
                                        break;
                                    }
                                    // ──── inventory 道具（新增）────
                                    case 'inventory': {
                                        const itemId = itemIds[idx % itemIds.length];
                                        if (invRound % 4 === 1) {
                                            ws.send(encodeMsg(nextSeq(), 'inventory.add', { itemId, count: 100 }));
                                            stats.sent++; stats.sentInventory++;
                                            stats.sentByProtocol['inventory.add']++;
                                        } else if (invRound % 4 === 2) {
                                            ws.send(encodeMsg(nextSeq(), 'inventory.consume', { itemId, count: 10 }));
                                            stats.sent++; stats.sentInventory++;
                                            stats.sentByProtocol['inventory.consume']++;
                                        } else if (invRound % 4 === 3) {
                                            ws.send(encodeMsg(nextSeq(), 'inventory.drop', { itemId, count: 5 }));
                                            stats.sent++; stats.sentInventory++;
                                            stats.sentByProtocol['inventory.drop']++;
                                        } else {
                                            ws.send(encodeMsg(nextSeq(), 'inventory.get', {}));
                                            stats.sent++; stats.sentInventory++;
                                            stats.sentByProtocol['inventory.get']++;
                                        }
                                        break;
                                    }
                                }
                            };
                            sendTimer = setInterval(sendOne, interval);
                        } else {
                            loginRetries++;
                            if (loginRetries === 1) {
                                ws.send(encodeMsg(nextSeq(), 'account.register', { username, password }));
                            } else if (loginRetries <= MAX_RETRIES) {
                                ws.send(encodeMsg(nextSeq(), 'account.login', { username, password }));
                            } else {
                                stats.loginFail++; loginDone = true;
                            }
                        }
                    } else if (data.name === 'account.register') {
                        ws.send(encodeMsg(nextSeq(), 'account.login', { username, password }));
                    }
                } catch {}
            }
        });

        ws.on('error', (e) => {
            stats.errors++;
            if (stats.errorSamples.length < 20) {
                stats.errorSamples.push({ idx, name: 'WS_ERROR', code: 'ERR', msg: (e && e.message) ? e.message.substring(0, 120) : 'unknown' });
            }
        });
        ws.on('close', (code, reason) => {
            stats.connected--;
            stats.disconnectCount++;
            const rc = String(code);
            stats.closeReasons[rc] = (stats.closeReasons[rc] || 0) + 1;
            if (sendTimer) clearInterval(sendTimer);
            // 只打印前 10 条 close 原因
            if (stats.disconnectCount <= 10) {
                const reasonStr = (reason || '').toString().substring(0, 60);
                const modeTag = sendMode ? `mode=${sendMode}` : '';
                console.log(`  [CLOSE #${stats.disconnectCount}] idx=${idx} code=${code} reason="${reasonStr}" loginOk=${loginDone} ${modeTag}`);
            }
        });
        setTimeout(() => { if (!ready) resolve(); }, 15000);
    });
}

function summary() {
    const elapsed = (performance.now() - stats.startAt) / 1000;
    const recvRatio = stats.sent > 0 ? ((stats.recv / stats.sent) * 100).toFixed(1) : 0;
    console.log(`\n═══════════════════════════════════════════════════════`);
    console.log(`📋 压测总结  mode=${args.mode}`);
    if (args.mode === 'mixed') {
        console.log(`   比例: world=${args.worldRatio} private=${args.privateRatio} room=${args.roomRatio} group=${args.groupRatio}`);
    } else if (args.mode === 'room') {
        console.log(`   roomSize=${args.roomSize}`);
    } else if (args.mode === 'group') {
        console.log(`   groupId=${args.groupId}`);
    } else if (args.mode === 'inventory') {
        console.log(`   itemIds=${args.itemIds}`);
    }
    console.log(`   总时长: ${elapsed.toFixed(1)}s  连接建立耗时: ${stats.connectTime}s`);
    console.log(`   连接数: ${args.conns}  活跃: ${stats.connected}`);
    console.log(`   Login: ✅${stats.loginOk}  ❌${stats.loginFail}`);
    console.log(`\n   【分模式发送】`);
    console.log(`     world:     ${stats.sentWorld}`);
    console.log(`     private:   ${stats.sentPrivate}`);
    console.log(`     room:      ${stats.sentRoom}`);
    console.log(`     group:     ${stats.sentGroup}`);
    console.log(`     inventory: ${stats.sentInventory}`);
    console.log(`\n   【分协议发送 Top 10】`);
    Object.entries(stats.sentByProtocol)
        .sort((a, b) => b[1] - a[1])
        .filter(([_, v]) => v > 0)
        .slice(0, 10)
        .forEach(([n, c]) => console.log(`     ${n.padEnd(22)} ${c}`));
    console.log(`\n   【汇总】`);
    console.log(`   发送: ${stats.sent}  QPS=${(stats.sent / elapsed).toFixed(0)}/s`);
    console.log(`   接收: ${stats.recv}  QPS=${(stats.recv / elapsed).toFixed(0)}/s  recv/send=${recvRatio}%`);
    console.log(`   业务错误: ${stats.errors}  断线: ${stats.disconnectCount}  最终活跃: ${stats.connected}`);

    // —— 错误分布 ——
    if (stats.errors > 0) {
        console.log(`\n   【分协议错误 Top 10】`);
        Object.entries(stats.errorsByProtocol)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 10)
            .forEach(([n, c]) => console.log(`     ${n.padEnd(22)} ${c}`));
        console.log(`\n   【错误码分布】`);
        Object.entries(stats.errorsByCode)
            .sort((a, b) => b[1] - a[1])
            .forEach(([c, n]) => console.log(`     code=${c.padEnd(6)} ${n}`));
    }

    // —— close 原因分布 ——
    if (stats.disconnectCount > 0) {
        console.log(`\n   【Close Code 分布】`);
        const closeCodeMap = {
            '1000': 'Normal Closure',
            '1001': 'Going Away',
            '1002': 'Protocol Error',
            '1003': 'Unsupported Data',
            '1006': 'Abnormal Closure',
            '1007': 'Invalid Frame Payload',
            '1008': 'Policy Violation',
            '1009': 'Message Too Big',
            '1010': 'Mandatory Extension',
            '1011': 'Internal Server Error',
            '1012': 'Service Restart',
            '1013': 'Try Again Later',
            '1014': 'Bad Gateway',
            '1015': 'TLS Handshake',
        };
        Object.entries(stats.closeReasons)
            .sort((a, b) => b[1] - a[1])
            .forEach(([code, n]) => {
                const label = closeCodeMap[code] || 'Unknown';
                console.log(`     ${code.padEnd(6)} ${label.padEnd(20)} ${n}`);
            });
    }

    // —— 错误样本 ——
    if (stats.errorSamples.length > 0) {
        console.log(`\n   【错误样本 (前 10 条)】`);
        stats.errorSamples.slice(0, 10).forEach((s, i) => {
            console.log(`     #${i + 1} idx=${s.idx} ${s.name} code=${s.code}  ${s.msg}`);
        });
    }

    console.log(`═══════════════════════════════════════════════════════\n`);
}

run().catch(e => { console.error('fatal:', e); process.exit(1); });