// bench.js — 压测 Broker 代理的并发能力
/**
   $env:NODE_PATH = (npm root -g)
   node .\testBench.js --conns=1000 --rate=10 --duration=30 --mode=mixed --worldRatio=0.1
**/

const WebSocket = require('ws');
const http = require('http');
const { performance } = require('perf_hooks');

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
            ws.send(JSON.stringify({ name: 'account.login', msg: { username, password } }));
            resolve();
        });

        ws.on('message', (data) => {
            const text = data.toString();
            stats.recv++;

            // 🟢 DEBUG: 前 10 条消息打印出来看后端返回什么
            if (stats.recv <= 10) {
                console.log(`  [DEBUG recv #${stats.recv}] ${text.substring(0, 300)}`);
            }

            if (loginDone) return;

            try {
                const msg = JSON.parse(text);
                if (msg.name === 'account.login') {
                    // 后端 WsMsgResp 结构: {seq, name, code, msg:{token,uid,...}}
                    if (msg.code === 0 && msg.msg && msg.msg.token) {
                        loginDone = true;
                        stats.loginOk++;
                        const tgt = (idx + 1) % args.conns;
                        sendTimer = setInterval(() => {
                            if (ws.readyState !== WebSocket.OPEN) return;
                            if (sendMode === 'world') {
                                ws.send(JSON.stringify({ name: 'chat.world.send', msg: { content: `bench_${idx}_${Date.now()%100000}` } }));
                                stats.sent++; stats.sentWorld++;
                            } else {
                                ws.send(JSON.stringify({ name: 'chat.private.send', msg: { toUid: `bench_${tgt}`, content: `hi_${idx}` } }));
                                stats.sent++; stats.sentPrivate++;
                            }
                        }, interval);
                    } else {
                        loginRetries++;
                        if (loginRetries === 1) {
                            ws.send(JSON.stringify({ name: 'account.register', msg: { username, password } }));
                        } else if (loginRetries <= MAX_RETRIES) {
                            ws.send(JSON.stringify({ name: 'account.login', msg: { username, password } }));
                        } else {
                            stats.loginFail++; loginDone = true;
                        }
                    }
                } else if (msg.name === 'account.register') {
                    ws.send(JSON.stringify({ name: 'account.login', msg: { username, password } }));
                }
            } catch {}
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