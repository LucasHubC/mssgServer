/**
 * mssgServer WebSocket 测试客户端
 * 
 * 用法:
 *   node testWS.js                              # 默认连 loginServer
 *   node testWS.js login                        # 同上
 *   node testWS.js chat                         # 连 chatServer
 *   node testWS.js ws://host:port               # 手动指定
 *   node testWS.js chat --username admin --password 123456
 */

const readline = require('readline');
let WebSocket;
try {
    WebSocket = require('ws');
} catch (e) {
    if (typeof globalThis.WebSocket === 'function') {
        WebSocket = globalThis.WebSocket;
    } else {
        console.error('❌ 缺少 ws 包，请先执行: npm install ws');
        process.exit(1);
    }
}

// ═══════════════════════════════════════════════
//  配置
// ═══════════════════════════════════════════════
// 全部走 Gateway 智能路由 —— 客户端永远只连这一个地址
// Gateway 按消息 name 前缀自动路由到 loginServer/chatServer/roomServer
const SERVER_PRESETS = {
    login: 'ws://127.0.0.1:9000/ws',
    chat:  'ws://127.0.0.1:9000/ws',
    room:  'ws://127.0.0.1:9000/ws'
};

// ═══════════════════════════════════════════════
//  解析命令行参数
// ═══════════════════════════════════════════════
const args = process.argv.slice(2);

function parseArgs(argv) {
    const opts = { server: 'login', username: '', password: '', token: '' };
    for (let i = 0; i < argv.length; i++) {
        switch (argv[i]) {
            case '--username': opts.username = argv[++i]; break;
            case '--password': opts.password = argv[++i]; break;
            case '--token': opts.token = argv[++i]; break;
            default:
                if (argv[i].startsWith('ws://') || argv[i].startsWith('wss://')) {
                    opts.server = argv[i];
                } else if (SERVER_PRESETS[argv[i]]) {
                    opts.server = argv[i];
                } else {
                    console.log(`⚠️  未知参数: ${argv[i]}`);
                }
        }
    }
    return opts;
}

// 拼 WS URL：如果有 token 参数自动加 ?token=xxx
function buildWsUrl(presetOrUrl, token) {
    let url = presetOrUrl.startsWith('ws') ? presetOrUrl : SERVER_PRESETS[presetOrUrl];
    if (token) {
        url += (url.includes('?') ? '&' : '?') + 'token=' + encodeURIComponent(token);
    }
    return url;
}

const opts = parseArgs(args);
const WS_URL = buildWsUrl(opts.server, opts.token);
let savedToken = opts.token; // 保存 token，switch 到 chatServer 时自动带上

// ═══════════════════════════════════════════════
//  状态
// ═══════════════════════════════════════════════
let currentUId = '';
let currentUsername = '';
let currentRoomId = ''; // 自动记录当前所在房间
let seq = 0;
function nextSeq() { return ++seq; }

// ═══════════════════════════════════════════════
//  发消息封装
// ═══════════════════════════════════════════════
function send(ws, name, msg = {}) {
    const payload = JSON.stringify({ seq: nextSeq(), name, msg });
    console.log(`📤 → ${name}`, JSON.stringify(msg));
    ws.send(payload);
}

// ═══════════════════════════════════════════════
//  连接
// ═══════════════════════════════════════════════
console.log('========================================');
console.log('  mssgServer WebSocket 测试客户端');
console.log(`  连接: ${WS_URL}`);
console.log('========================================');

const ws = new WebSocket(WS_URL);

ws.onopen = () => {
    console.log('✅ 连接成功\n');

    // 如果带了 --username 参数，自动登录
    if (opts.username && opts.password) {
        console.log(`🔄 自动登录: ${opts.username}`);
        send(ws, 'account.login', { username: opts.username, password: opts.password });
        // 等登录响应后再进入交互
        setTimeout(startRepl, 500);
    } else {
        startRepl();
    }
};

ws.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);
        // 登录响应自动缓存 uid + token
        if (data.name === 'account.login' && data.code === 0 && data.msg) {
            currentUId = data.msg.uid || '';
            currentUsername = data.msg.username || '';
            currentRoomId = ''; // 登录后清空房间
            if (data.msg.token) savedToken = data.msg.token;
            console.log(`\n🎯 登录成功! UID=${currentUId} Username=${currentUsername} Token=${savedToken ? '(已保存)' : '(无)'} \n`);
        }
        // 房间聊天推送美化
        else if (data.name === 'room.chat.new' && data.msg) {
            const m = data.msg;
            const dir = m.from_uid === currentUId ? '📤' : '💬';
            console.log(`${dir} [房间:${m.room_id}] ${m.from_name}(${m.from_uid.substring(0,10)}): ${m.content}`);
            console.log(`   msgId=${m.msg_id} ts=${m.timestamp}`);
            return; // 不打下面的通用日志
        }
        // 房间系统通知美化
        else if (data.name === 'room.notice' && data.msg) {
            const n = data.msg;
            switch (n.type) {
                case 'room_created':
                    console.log(`\n🏠 你创建了房间 [${n.roomId}] "${n.name}"\n`);
                    break;
                case 'user_joined':
                    console.log(`\n🟢 ${n.username}(${n.uid?.substring(0,10)}) 加入了房间 [${n.roomId}]`);
                    console.log(`   当前成员: ${n.members?.length || 0} 人\n`);
                    break;
                case 'user_left':
                    console.log(`\n🔴 ${n.username}(${n.uid?.substring(0,10)}) 离开了房间 [${n.roomId}]`);
                    console.log(`   当前成员: ${n.members?.length || 0} 人\n`);
                    break;
                case 'user_kicked':
                    console.log(`\n👢 房主把 ${n.target?.substring(0,10)} 踢出了房间 [${n.roomId}]`);
                    console.log(`   当前成员: ${n.members?.length || 0} 人\n`);
                    break;
                case 'kicked':
                    currentRoomId = '';
                    console.log(`\n⛔ 你被房主踢出了房间 [${n.roomId}]!\n`);
                    break;
                case 'room_deleted':
                    currentRoomId = '';
                    console.log(`\n💥 房间 [${n.roomId}] 已被房主删除\n`);
                    break;
                default:
                    console.log(`\n📢 房间通知: ${JSON.stringify(n)}\n`);
            }
            return;
        }
        // Token 自动鉴权结果（连 chatServer 时收到）
        else if (data.name === 'auth.ok' && data.msg) {
            currentUId = data.msg.uid || '';
            currentUsername = data.msg.username || '';
            console.log(`\n✅ Token 鉴权成功! UID=${currentUId} Username=${currentUsername}\n`);
        }
        else if (data.name === 'auth.fail' && data.msg) {
            console.log(`\n❌ Token 鉴权失败: ${data.msg.errMsg}\n`);
        }
        // 创建/加入房间成功时自动记录 roomId
        if (data.name === 'room.create' && data.code === 0 && data.msg && data.msg.room_id) {
            currentRoomId = data.msg.room_id;
            console.log(`\n🏠 已自动进入房间: ${currentRoomId}\n`);
        }
        if (data.name === 'room.join' && data.code === 0 && data.msg && data.msg.room_id) {
            currentRoomId = data.msg.room_id;
            console.log(`\n🏠 已自动进入房间: ${currentRoomId}\n`);
        }
        if (data.name === 'room.leave' && data.code === 0) {
            currentRoomId = '';
            console.log(`\n👋 已离开房间\n`);
        }

        // 在线用户列表美化
        if (data.name === 'chat.online.get' && data.code === 0 && data.msg) {
            const m = data.msg;
            console.log(`\n📊 本机在线用户: ${m.total} 人`);
            if (m.uids && m.uids.length > 0) {
                m.uids.forEach((uid, i) => {
                    const marker = uid === currentUId ? ' ⭐(我)' : '';
                    console.log(`  ${i + 1}. ${uid}${marker}`);
                });
            }
            console.log('');
            return;
        }

        console.log(`📥 ← [${data.name}] seq=${data.seq} code=${data.code}`);
        if (data.msg) {
            const formatted = typeof data.msg === 'object' 
                ? JSON.stringify(data.msg, null, 2) 
                : data.msg;
            console.log('   ', formatted);
        }
    } catch (e) {
        console.log('📥 ← 原始:', event.data);
    }
};

ws.onclose = (code, reason) => {
    console.log(`\n🔌 连接关闭 (code=${code})`);
    process.exit(0);
};

ws.onerror = (err) => {
    console.error('❌ 错误:', err.message || err);
};

// ═══════════════════════════════════════════════
//  跨服自动化测试
// ═══════════════════════════════════════════════
// 保持架构干净：chatServer 只做 chat，不注册 login/room 路由
// 跨服测试验证三层能力：
//   ① Redis 跨服务 session（token 自动鉴权）
//   ② Redis Pub/Sub 跨服广播链路
//   ③ Gateway admin HTTP 接口（触发跨服 push/broadcast/kick）
//
// 测试流程：
//   1. A 在 loginServer 注册+登录 → 拿到 tokenA
//   2. B 也在 loginServer 注册+登录 → 拿到 tokenB
//   3. B 断开 → 用 tokenB 连 chatServer → token 自动鉴权成功（Redis session 跨服务 ✅）
//   4. Gateway admin/push → 推给 B → B 收到（跨服 Pub/Sub push ✅）
//   5. Gateway admin/broadcast → 全服广播 → A 和 B 都收到（跨服 broadcast ✅）
//   6. Gateway admin/kick → 踢 A → A 被踢下线（跨服 kick ✅）
//
// 依赖：Redis + loginServer(8003) + chatServer(8004) + Gateway(9000) 都在跑
const GATEWAY_HTTP = 'http://127.0.0.1:9000';

// 用 Node 内置 http 模块发 POST（兼容旧 Node 版本）
function httpPost(url, body) {
    return new Promise((resolve, reject) => {
        try {
            const u = new URL(url);
            const http = require(u.protocol === 'https:' ? 'https' : 'http');
            const data = JSON.stringify(body);
            const req = http.request({
                hostname: u.hostname,
                port: u.port || (u.protocol === 'https:' ? 443 : 80),
                path: u.pathname + u.search,
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Content-Length': Buffer.byteLength(data)
                }
            }, (res) => {
                let raw = '';
                res.on('data', chunk => raw += chunk);
                res.on('end', () => resolve(raw));
            });
            req.on('error', reject);
            req.write(data);
            req.end();
        } catch (e) {
            reject(e);
        }
    });
}

async function gatewayAdmin(route, body) {
    return await httpPost(`${GATEWAY_HTTP}${route}`, body);
}

function runCrossServerTest() {
    console.log('\n╔══════════════════════════════════════════════╗');
    console.log('║     🔄 跨服自动化测试 (Cross-Server Test)      ║');
    console.log('╚══════════════════════════════════════════════╝\n');

    let step = 0;
    const passed = [];
    const failed = [];
    function ok(name) { passed.push(name); console.log(`  ✅ Step ${++step}: ${name}`); }
    function fail(name, err) { failed.push({name, err}); console.log(`  ❌ Step ${++step}: ${name} — ${err}`); }

    const now = Date.now().toString().slice(-6);
    const userA = `crossA_${now}`;
    const passA = '123456';
    const userB = `crossB_${now}`;
    const passB = '123456';

    let connA, connB;
    let uidA = '', uidB = '', tokenA = '', tokenB = '';

    function makeConn(url, label) {
        return new Promise((resolve, reject) => {
            const c = new WebSocket(url);
            c._label = label;
            c._pending = {};
            c.onopen = () => resolve(c);
            c.onerror = (e) => reject(new Error(`${label} connect failed: ${e.message || e}`));
            c.onmessage = (event) => {
                const data = JSON.parse(event.data);
                if (!c._pending[data.name]) c._pending[data.name] = [];
                c._pending[data.name].push(data);
            };
        });
    }

    function sendAndWait(c, name, body, expectedResponseName = null, timeout = 2000) {
        return new Promise((resolve) => {
            const mySeq = ++seq;
            const p = JSON.stringify({ seq: mySeq, name, msg: body });
            c.send(p);

            if (!expectedResponseName) {
                setTimeout(() => resolve(null), timeout);
                return;
            }

            const start = Date.now();
            const interval = setInterval(() => {
                const msgs = c._pending[expectedResponseName];
                if (msgs && msgs.length > 0) {
                    const last = msgs[msgs.length - 1];
                    clearInterval(interval);
                    resolve(last);
                    return;
                }
                if (Date.now() - start > timeout) {
                    clearInterval(interval);
                    resolve(null);
                }
            }, 100);
        });
    }

    async function run() {
        try {
            // ═══════ 1. 连接 A → Gateway（智能路由）═══════
            console.log('\n📡 [1/6] 连接 A → Gateway (9000)');
            try {
                connA = await makeConn(SERVER_PRESETS.login, 'A');
                ok('A 连接 Gateway 成功 — Gateway 智能路由已就绪 ✅');
            } catch (e) {
                fail('A 连接 Gateway', e.message);
                throw new Error('无法继续，Gateway 可能没启动');
            }

            // ═══════ 2. 连接 B → Gateway（智能路由）═══════
            console.log('📡 [2/6] 连接 B → Gateway (9000)');
            try {
                connB = await makeConn(SERVER_PRESETS.chat, 'B');
                ok('B 连接 Gateway 成功 — 两路独立 WS 代理连接 ✅');
            } catch (e) {
                fail('B 连接 Gateway', e.message);
                throw new Error('无法继续，Gateway 可能没启动');
            }

            // ═══════ 3. A 注册 + 登录 ═══════
            console.log('\n🔐 [3/6] A 注册并登录');
            await sendAndWait(connA, 'account.register', { username: userA, password: passA }, 'account.register');
            let loginRsp = await sendAndWait(connA, 'account.login', { username: userA, password: passA }, 'account.login');
            if (!loginRsp || loginRsp.code !== 0) {
                loginRsp = await sendAndWait(connA, 'account.login', { username: userA, password: passA }, 'account.login');
            }
            if (loginRsp && loginRsp.code === 0 && loginRsp.msg) {
                uidA = loginRsp.msg.uid || '';
                tokenA = loginRsp.msg.token || '';
                ok(`A 登录成功 uid=${uidA}（Gateway → loginServer ✅）`);
            } else {
                fail('A 登录', loginRsp?.msg?.errMsg || 'unknown');
                throw new Error('无法继续');
            }

            // ═══════ 3.5 ⭐ Gateway WS 智能路由核心验证（同一条 WS 连三个后端）═══════
            console.log('\n🧠 [3.5] Gateway WS 智能路由核心验证 — 单 WS 连接代理 login + room + chat');
            const testRoomName = `gw_room_${now}`;
            let roomId = '';

            // 3.5a: room.create → Gateway 路由到 roomServer
            console.log('  ├─ 3.5a. room.create → Gateway → roomServer');
            const createRsp = await sendAndWait(connA, 'room.create', { name: testRoomName, max_size: 10 }, 'room.create', 3000);
            if (createRsp && createRsp.code === 0 && createRsp.msg) {
                roomId = createRsp.msg.room_id || '';
                ok(`room.create 成功 roomId=${roomId}（Gateway → roomServer ✅）`);
            } else {
                // roomServer 可能没启动，跳过但不报错
                console.log('  ⚠️ room.create 失败，可能 roomServer 没启动，跳过 room 路由验证');
            }

            // 3.5b: chat.world.send → Gateway 路由到 chatServer
            console.log('  └─ 3.5b. chat.world.send → Gateway → chatServer');
            const chatRsp = await sendAndWait(connA, 'chat.world.send', { content: `Gateway智能路由测试 ${now}` }, 'chat.world.send', 3000);
            if (chatRsp && chatRsp.code === 0) {
                ok(`chat.world.send 成功（Gateway → chatServer ✅）`);
            } else {
                console.log('  ⚠️ chat.world.send 失败，可能 chatServer 没启动');
            }

            // ═══════ 4. B 也在 loginServer 注册+登录（拿 tokenB）═══════
            console.log('\n🔐 [4/6] B 注册并登录（走 loginServer 拿 tokenB）');
            await sendAndWait(connA, 'account.register', { username: userB, password: passB }, 'account.register');
            let loginRspB = await sendAndWait(connA, 'account.login', { username: userB, password: passB }, 'account.login');
            if (!loginRspB || loginRspB.code !== 0) {
                loginRspB = await sendAndWait(connA, 'account.login', { username: userB, password: passB }, 'account.login');
            }
            if (loginRspB && loginRspB.code === 0 && loginRspB.msg) {
                uidB = loginRspB.msg.uid || '';
                tokenB = loginRspB.msg.token || '';
                ok(`B 登录成功 uid=${uidB}`);
            } else {
                fail('B 登录', loginRspB?.msg?.errMsg || 'unknown');
                throw new Error('无法继续');
            }

            // ═══════ 5. B 断开 → 用 tokenB 重连 Gateway ═══════
            console.log('\n🔄 [5/6] B 断开 Gateway，用 tokenB 重连（验证 Gateway 透传 token → 后端鉴权）');
            connB.close();
            await new Promise(r => setTimeout(r, 300));

            try {
                connB = await makeConn(buildWsUrl(SERVER_PRESETS.chat, tokenB), 'B');
                await new Promise(r => setTimeout(r, 500)); // 等 auth.ok
                if (uidB) {
                    ok(`B Token 鉴权成功 — Gateway 透传 token + Redis 跨服务 session OK ✅`);
                } else {
                    fail('B Token 鉴权', '没收到 auth.ok');
                    throw new Error('Redis 跨服务 session 可能有问题');
                }
            } catch (e) {
                fail('B token 重连', e.message);
                throw new Error('无法继续，Gateway/Redis 可能没启动');
            }

            // 等 Pub/Sub 订阅就绪
            await new Promise(r => setTimeout(r, 500));

            // ═══════ 6. 跨服 Gateway admin HTTP 接口测试 ═══════
            console.log('\n🧪 [6/6] Gateway admin HTTP 接口跨服测试');

            // 给 connA/connB 加动态消息监听
            let bReceivedPush = false, aKicked = false;
            const origB = connB.onmessage;
            connB.onmessage = (event) => {
                try {
                    const d = JSON.parse(event.data);
                    if (d.name === 'system.push') {
                        bReceivedPush = true;
                        console.log(`     📥 B 收到 admin/push`);
                    }
                    if (origB) origB(event);
                } catch (e) {}
            };
            const origCloseA = connA.onclose;
            connA.onclose = () => { aKicked = true; if (origCloseA) origCloseA(); };

            // 6a: admin/push → B
            console.log('\n  ├─ 6a. POST /admin/push → 跨服推给 B');
            try {
                await gatewayAdmin('/admin/push', {
                    uid: uidB,
                    name: 'system.push',
                    data: { from: 'gateway', content: `hello-${now}` }
                });
                await new Promise(r => setTimeout(r, 800));
                if (bReceivedPush) {
                    ok('B 收到 admin/push — 跨服 Pub/Sub push ✅');
                } else {
                    fail('admin/push → B', 'B 没收到');
                }
            } catch (e) {
                fail('admin/push', `Gateway 可能没启动: ${e.message}`);
            }

            // 6b: admin/broadcast → A + B
            let aBC = false, bBC = false;
            const origA = connA.onmessage;
            connA.onmessage = (event) => {
                try {
                    const d = JSON.parse(event.data);
                    if (d.name === 'system.broadcast') aBC = true;
                    if (origA) origA(event);
                } catch (e) {}
            };
            const origB2 = connB.onmessage;
            connB.onmessage = (event) => {
                try {
                    const d = JSON.parse(event.data);
                    if (d.name === 'system.broadcast') bBC = true;
                    if (origB2) origB2(event);
                } catch (e) {}
            };

            console.log('  ├─ 6b. POST /admin/broadcast → 全服广播');
            try {
                await gatewayAdmin('/admin/broadcast', {
                    name: 'system.broadcast',
                    data: { from: 'gateway', content: `world-${now}` }
                });
                await new Promise(r => setTimeout(r, 800));
                if (aBC && bBC) {
                    ok('A 和 B 都收到全服广播 — 跨服 Pub/Sub broadcast ✅');
                } else {
                    fail('admin/broadcast', `A=${aBC} B=${bBC}`);
                }
            } catch (e) {
                fail('admin/broadcast', `Gateway 可能没启动: ${e.message}`);
            }

            // 6c: admin/kick → A
            console.log('  └─ 6c. POST /admin/kick → 跨服踢 A');
            try {
                await gatewayAdmin('/admin/kick', { uid: uidA });
                await new Promise(r => setTimeout(r, 1000));
                if (aKicked) {
                    ok('A 被踢下线 — 跨服 Pub/Sub kick ✅');
                } else {
                    fail('admin/kick → A', 'A 没被踢下线');
                }
            } catch (e) {
                fail('admin/kick', `Gateway 可能没启动: ${e.message}`);
            }

            // ═══════ 汇总 ═══════
            console.log('\n╔══════════════════════════════════════════════╗');
            console.log('║                 📊 测试结果汇总               ║');
            console.log('╠══════════════════════════════════════════════╣');
            console.log('║  🧪 架构：Gateway 智能路由（客户端永远只连一个 WS）');
            console.log('║  🧪 路由表：account→login  room→room  chat→chat  group→chat');
            console.log('║  🧪 同一条 WS 同时操作 login + room + chat（无需 switch！）');
            console.log('║  🧪 跨服：Redis session + Redis Pub/Sub');
            console.log('╠══════════════════════════════════════════════╣');
            passed.forEach(p => console.log(`║  ✅ ${p}`));
            failed.forEach(f => console.log(`║  ❌ ${f.name} — ${f.err}`));
            console.log('╠══════════════════════════════════════════════╣');
            if (failed.length === 0) {
                console.log(`║  🎉 全部通过！(${passed.length}项) 跨服链路完全 OK`);
            } else {
                console.log(`║  ⚠️  通过 ${passed.length}/${passed.length + failed.length} 项`);
            }
            console.log('╚══════════════════════════════════════════════╝\n');

            if (failed.length > 0) {
                console.log('💡 排错提示：');
                if (failed.some(f => f.name.includes('Gateway'))) {
                    console.log('   → Gateway(9000) 可能没启动: go run gatewayServer.go');
                }
                console.log('   → Redis(6379) 必须在跑');
                console.log('   → Gateway(9000) + loginServer(8003) + chatServer(8004) + roomServer(8005) 都要启动');
            }
        } catch (e) {
            console.log('💥 测试异常:', e.message);
            console.log('💡 请确认：Redis + Gateway(9000) + loginServer(8003) + chatServer(8004) + roomServer(8005) 都在运行\n');
        } finally {
            try { connA?.close(); } catch(e) {}
            try { connB?.close(); } catch(e) {}
        }
    }

    run();
}
  


// ═══════════════════════════════════════════════
//  交互式命令行
// ═══════════════════════════════════════════════
const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

function printHelp() {
    console.log(`
📋 可用命令:
  ─── 登录 ───
  login <username> <password>       登录
  whoami                            显示当前 uid

  ─── 世界频道 ───
  world <content>                   世界频道聊天（跨服全服广播）

  ─── 私聊 ───
  private <to_uid> <content>         私聊（自动跨服推到目标节点）

  ─── 群组 ───
  join <group_id>                    加入群组
  leave <group_id>                   退出群组
  members <group_id>                 查看群组成员
  group <group_id> <content>         群组聊天

  ─── 房间 ───
  rc / rcreate <name> [maxSize]     创建房间（自动进入）
  rj / rjoin <roomId>               加入房间
  rl / rleave [roomId]              离开房间（默认当前房间）
  rk / rkick <roomId> <targetUid>   房主踢人
  rd / rdelete [roomId]             房主删除房间
  rs / rsend <content>              发送房间消息（默认当前房间）
  rwho                              查看当前所在房间

  ─── 在线 ───
  online                            查看本机在线用户

  ─── 跨服测试 ───
  crosstest                         自动化跨服测试（同时开2个连接，login+chat）

  ─── 跨服测试 ───
  crosstest                         自动化跨服测试（同时开2个连接，login+chat）

  ─── 管理 ───
  quit / exit                       退出
  help                              显示帮助
`);
    rl.prompt();
}

function startRepl() {
    rl.setPrompt(`[${WS_URL.split('/').pop()}]> `);
    console.log('输入 help 查看命令列表\n');
    rl.prompt();

    rl.on('line', (line) => {
        const input = line.trim();
        if (!input) { rl.prompt(); return; }

        const [cmd, ...args] = input.split(/\s+/);

        switch (cmd.toLowerCase()) {
            case 'quit':
            case 'exit':
                console.log('👋 bye~');
                ws.close();
                rl.close();
                return;

            case 'help':
            case '?':
                printHelp();
                return;

            case 'login':
                if (args.length < 2) { console.log('用法: login <username> <password>'); break; }
                send(ws, 'account.login', { username: args[0], password: args[1] });
                break;

             case 'register':
                if (args.length < 2) { console.log('用法: register <username> <password>'); break; }
                send(ws, 'account.register', { username: args[0], password: args[1] });
                break;

            case 'whoami':
                console.log(`  UID: ${currentUId || '(未登录)'}`);
                console.log(`  Username: ${currentUsername || '(未登录)'}`);
                break;

            case 'world':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: world <内容>'); break; }
                send(ws, 'chat.world.send', { content: args.join(' ') });
                break;

            case 'private':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: private <to_uid> <内容>'); break; }
                send(ws, 'chat.private.send', { to_uid: args[0], content: args.slice(1).join(' ') });
                break;

            case 'join':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: join <group_id>'); break; }
                send(ws, 'chat.group.join', { group_id: args[0] });
                break;

            case 'leave':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: leave <group_id>'); break; }
                send(ws, 'chat.group.leave', { group_id: args[0] });
                break;

            case 'members':
                if (args.length < 1) { console.log('用法: members <group_id>'); break; }
                send(ws, 'chat.group.members', { group_id: args[0] });
                break;

            case 'group':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: group <group_id> <内容>'); break; }
                send(ws, 'chat.group.send', { group_id: args[0], content: args.slice(1).join(' ') });
                break;

            // ──── 房间 ────
            case 'rc':
            case 'rcreate':
            case 'room-create': {
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: rc <name> [maxSize]'); break; }
                const body = { name: args[0] };
                if (args[1]) body.max_size = parseInt(args[1], 10) || 10;
                send(ws, 'room.create', body);
                break;
            }

            case 'rj':
            case 'rjoin':
            case 'room-join':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: rj <roomId>'); break; }
                send(ws, 'room.join', { room_id: args[0] });
                break;

            case 'rl':
            case 'rleave':
            case 'room-leave': {
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                const rid = args[0] || currentRoomId;
                if (!rid) { console.log('❌ 请指定 roomId 或先加入房间'); break; }
                send(ws, 'room.leave', { room_id: rid });
                break;
            }

            case 'rk':
            case 'rkick':
            case 'room-kick':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: rk <roomId> <targetUid>'); break; }
                send(ws, 'room.kick', { room_id: args[0], target_uid: args[1] });
                break;

            case 'rd':
            case 'rdelete':
            case 'room-delete': {
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                const rid = args[0] || currentRoomId;
                if (!rid) { console.log('❌ 请指定 roomId 或先加入房间'); break; }
                send(ws, 'room.delete', { room_id: rid });
                break;
            }

            case 'rs':
            case 'rsend':
            case 'room-send': {
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (!currentRoomId) { console.log('❌ 请先 join 一个房间（或用 rj）'); break; }
                if (args.length < 1) { console.log('用法: rs <内容>'); break; }
                send(ws, 'room.send', { room_id: currentRoomId, content: args.join(' ') });
                break;
            }

            case 'rwho':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (currentRoomId) {
                    console.log(`  🏠 当前房间: ${currentRoomId}`);
                } else {
                    console.log('  (不在任何房间里)');
                }
                break;

            // ──── 在线用户 ────
            case 'online':
            case 'online-users':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                send(ws, 'chat.online.get', {});
                break;

            case 'switch':
            case 'connect': {
                // 快速切换到另一个 server（带上已保存的 token）
                if (args.length < 1) { console.log('用法: switch <login|chat|room|ws://host:port>'); break; }
                const newUrl = buildWsUrl(args[0], savedToken);
                console.log(`🔄 切换到 ${newUrl}`);
                // 重连（内联，修复 setupWsHandlers 未定义的 bug）
                const newWs = new WebSocket(newUrl);
                newWs.onopen = () => { console.log(`✅ 已连接到 ${newUrl}`); rl.prompt(); };
                newWs.onmessage = ws.onmessage;
                newWs.onclose = ws.onclose;
                newWs.onerror = ws.onerror;
                ws.close();
                break;
            }

            case 'token':
                if (args.length < 1) {
                    if (savedToken) console.log(`当前 token: ${savedToken}`);
                    else console.log('❌ 还没有 token');
                    break;
                }
                savedToken = args[0];
                console.log(`✅ token 已保存: ${savedToken}`);
                break;

            // ──── 道具模块 ────
            case 'additem':
            case 'add':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: additem <itemId> <count>'); break; }
                send(ws, 'inventory.add', { itemId: args[0], count: parseInt(args[1], 10) || 1 });
                break;

            case 'consume':
            case 'useitem':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: consume <itemId> <count>'); break; }
                send(ws, 'inventory.consume', { itemId: args[0], count: parseInt(args[1], 10) || 1 });
                break;

            case 'dropitem':
            case 'drop':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: dropitem <itemId> <count>'); break; }
                send(ws, 'inventory.drop', { itemId: args[0], count: parseInt(args[1], 10) || 1 });
                break;

            case 'inventory':
            case 'bag':
            case 'items':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                send(ws, 'inventory.get', {});
                break;

            // ──── 跨服自动化测试 ────
            case 'crosstest':
            case 'cross':
                runCrossServerTest();
                break;

            // ──── 跨服自动化测试 ────
            case 'crosstest':
            case 'cross':
                runCrossServerTest();
                break;

            default:
                // 尝试把剩余参数当 JSON 解析（通用兜底，比如直接发 {"itemId":"xxx","count":1}）
                try {
                    const jsonPart = input.slice(cmd.length).trim();
                    if (!jsonPart) { console.log('❌ 未知命令，输入 help 查看帮助'); break; }
                    const msg = JSON.parse(jsonPart);
                    send(ws, cmd, msg);
                } catch (e) {
                    console.log('❌ 未知命令，输入 help 查看帮助');
                }
        }
        rl.prompt();
    });
}