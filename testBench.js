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
    mode: 'mixed', worldRatio: 0.1
};
process.argv.slice(2).forEach(a => {
    const [k, v] = a.split('=');
    const key = k.replace(/^--/, '');
    if (args[key] !== undefined) args[key] = isNaN(v) ? v : +v;
});

const stats = {
    connected: 0,
    sent: 0, sentWorld: 0, sentPrivate: 0,
    recv: 0,
    errors: 0,
    loginOk: 0, loginFail: 0,
    connectTime: 0,
    startAt: 0,
};

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

        let sendMode = args.mode;
        if (args.mode === 'mixed') sendMode = Math.random() < args.worldRatio ? 'world' : 'private';

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

                if (loginDone) continue;

                try {
                    if (data.name === 'account.login') {
                        if (data.code === 0 && data.msg && data.msg.token) {
                            loginDone = true;
                            stats.loginOk++;
                            const tgt = (idx + 1) % args.conns;
                            sendTimer = setInterval(() => {
                                if (ws.readyState !== WebSocket.OPEN) return;
                                if (sendMode === 'world') {
                                    ws.send(encodeMsg(nextSeq(), 'chat.world.send', { content: `bench_${idx}_${Date.now()%100000}` }));
                                    stats.sent++; stats.sentWorld++;
                                } else {
                                    ws.send(encodeMsg(nextSeq(), 'chat.private.send', { toUid: `bench_${tgt}`, content: `hi_${idx}` }));
                                    stats.sent++; stats.sentPrivate++;
                                }
                            }, interval);
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

        ws.on('error', () => { stats.errors++; });
        ws.on('close', () => { stats.connected--; if (sendTimer) clearInterval(sendTimer); });
        setTimeout(() => { if (!ready) resolve(); }, 15000);
    });
}

function summary() {
    const elapsed = (performance.now() - stats.startAt) / 1000;
    const recvRatio = stats.sent > 0 ? ((stats.recv / stats.sent) * 100).toFixed(1) : 0;
    console.log(`\n═══════════════════════════════════════════`);
    console.log(`📋 压测总结  mode=${args.mode}  worldRatio=${args.worldRatio}`);
    console.log(`   总时长: ${elapsed.toFixed(1)}s`);
    console.log(`   连接数: ${args.conns}  活跃: ${stats.connected}`);
    console.log(`   Login: ✅${stats.loginOk}  ❌${stats.loginFail}`);
    console.log(`   发送: ${stats.sent} (world=${stats.sentWorld} private=${stats.sentPrivate})`);
    console.log(`   接收: ${stats.recv}  recv/send = ${recvRatio}%`);
    console.log(`   发送 QPS: ${(stats.sent / elapsed).toFixed(0)}/s`);
    console.log(`   接收 QPS: ${(stats.recv / elapsed).toFixed(0)}/s`);
    console.log(`   错误: ${stats.errors}`);
    console.log(`═══════════════════════════════════════════\n`);
}

run().catch(e => { console.error('fatal:', e); process.exit(1); });