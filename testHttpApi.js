/**
 * mssgServer HTTP API 测试客户端（短连接模式）
 * 
 * 用法:
 *   node testAPI.js                              # 默认连 http://127.0.0.1:8005
 *   node testAPI.js http://127.0.0.1:8005        # 手动指定
 *   node testAPI.js --username admin --password 123456   # 自动登录
 * 
 * 依赖: 零依赖！只用 Node.js 内置的 http/https + readline
 */

const http = require('http');
const https = require('https');
const readline = require('readline');
const { URL } = require('url');

// ═══════════════════════════════════════════════
//  解析命令行参数
// ═══════════════════════════════════════════════
const args = process.argv.slice(2);
function parseArgs(argv) {
    const opts = { baseUrl: 'http://127.0.0.1:8005', username: '', password: '' };
    for (let i = 0; i < argv.length; i++) {
        switch (argv[i]) {
            case '--username': opts.username = argv[++i]; break;
            case '--password': opts.password = argv[++i]; break;
            default:
                if (argv[i].startsWith('http')) opts.baseUrl = argv[i];
        }
    }
    // 去掉末尾的 /
    opts.baseUrl = opts.baseUrl.replace(/\/+$/, '');
    return opts;
}
const opts = parseArgs(args);

// ═══════════════════════════════════════════════
//  状态
// ═══════════════════════════════════════════════
let currentToken = '';
let currentUId = '';
let currentUsername = '';

// ═══════════════════════════════════════════════
//  HTTP 请求封装（零依赖，支持 Bearer Token）
// ═══════════════════════════════════════════════
function request(method, path, body = null, extraHeaders = {}) {
    return new Promise((resolve, reject) => {
        const fullUrl = opts.baseUrl + path;
        const urlObj = new URL(fullUrl);
        const lib = urlObj.protocol === 'https:' ? https : http;

        const headers = {
            'Content-Type': 'application/json',
            ...extraHeaders,
        };
        const bodyStr = body ? JSON.stringify(body) : null;
        if (bodyStr) headers['Content-Length'] = Buffer.byteLength(bodyStr);

        const req = lib.request({
            method,
            hostname: urlObj.hostname,
            port: urlObj.port || (urlObj.protocol === 'https:' ? 443 : 80),
            path: urlObj.pathname + urlObj.search,
            headers,
        }, (res) => {
            let data = '';
            res.on('data', (chunk) => data += chunk);
            res.on('end', () => {
                try {
                    const json = JSON.parse(data);
                    resolve({ status: res.statusCode, body: json });
                } catch (e) {
                    resolve({ status: res.statusCode, body: data });
                }
            });
        });

        req.on('error', reject);
        if (bodyStr) req.write(bodyStr);
        req.end();
    });
}

// 带 Token 的快捷方法
function authHeaders() {
    return currentToken ? { Authorization: `Bearer ${currentToken}` } : {};
}

async function get(path)    { return request('GET',    path, null, authHeaders()); }
async function post(path, body) { return request('POST', path, body, authHeaders()); }

// ═══════════════════════════════════════════════
//  打印工具
// ═══════════════════════════════════════════════
function printResult(resp, label = '') {
    const icon = resp.body.code === 0 ? '✅' : '❌';
    console.log(`${icon} ${label}`);
    console.log(`   HTTP ${resp.status}`);
    if (typeof resp.body === 'string') {
        console.log('   ', resp.body);
    } else {
        console.log('   ', JSON.stringify(resp.body, null, 2));
    }
}

// ═══════════════════════════════════════════════
//  命令处理
// ═══════════════════════════════════════════════
async function handleRegister(args) {
    if (args.length < 2) { console.log('用法: register <username> <password>'); return; }
    const resp = await post('/api/register', { username: args[0], password: args[1] });
    printResult(resp, '注册');
}

async function handleLogin(args) {
    if (args.length < 2) { console.log('用法: login <username> <password>'); return; }
    const resp = await post('/api/login', { username: args[0], password: args[1] });
    printResult(resp, '登录');
    if (resp.body.code === 0 && resp.body.data) {
        currentToken = resp.body.data.token;
        currentUId = resp.body.data.uid;
        currentUsername = resp.body.data.username;
        console.log(`   🎯 Token 已缓存: ${currentToken.substring(0, 20)}...`);
    }
}

async function handleLogout() {
    if (!currentToken) { console.log('❌ 未登录'); return; }
    const resp = await post('/api/logout', {});
    printResult(resp, '登出');
    if (resp.body.code === 0) {
        currentToken = '';
        currentUId = '';
        currentUsername = '';
        console.log('   🎯 本地 Token 已清除');
    }
}

async function handleMe() {
    if (!currentToken) { console.log('❌ 未登录，请先 login'); return; }
    const resp = await get('/api/me');
    printResult(resp, '当前用户信息 (/api/me)');
}

async function handleUserInfo(args) {
    if (args.length < 1) { console.log('用法: user <uid>'); return; }
    const resp = await get(`/api/user/${args[0]}`);
    printResult(resp, '查询用户');
}

async function handleWhoami() {
    console.log(`  Token:     ${currentToken ? currentToken.substring(0, 25) + '...' : '(无)'}`);
    console.log(`  UID:       ${currentUId || '(未登录)'}`);
    console.log(`  Username:   ${currentUsername || '(未登录)'}`);
    console.log(`  Token 有效: ${currentToken ? '是（Redis 里应该能搜到）' : '否'}`);
}

async function handleHealth() {
    try {
        const resp = await get('/');
        printResult(resp, '健康检查 GET /');
    } catch (e) {
        console.log('❌ 连接失败:', e.message);
    }
}

// ────────────────────────────────────────
//  道具模块
// ────────────────────────────────────────

async function handleItemAdd(args) {
    if (!currentToken) { console.log('❌ 未登录，请先 login'); return; }
    if (args.length < 2) { console.log('用法: additem <itemId> <count>'); return; }
    const count = parseInt(args[1], 10) || 1;
    const resp = await post('/api/inventory/add', { itemId: args[0], count });
    printResult(resp, `增加道具 ${args[0]} × ${count}`);
}

async function handleItemConsume(args) {
    if (!currentToken) { console.log('❌ 未登录，请先 login'); return; }
    if (args.length < 2) { console.log('用法: consume <itemId> <count>'); return; }
    const count = parseInt(args[1], 10) || 1;
    const resp = await post('/api/inventory/consume', { itemId: args[0], count });
    printResult(resp, `消耗道具 ${args[0]} × ${count}`);
}

async function handleItemGet() {
    if (!currentToken) { console.log('❌ 未登录，请先 login'); return; }
    const resp = await get('/api/inventory/get');
    // 美化输出背包
    if (resp.body.code === 0 && resp.body.data) {
        const items = resp.body.data.items || [];
        console.log(`
🎒 背包 (${items.length} 种道具):`);
        if (items.length === 0) {
            console.log('   (空)');
        } else {
            // 对齐列宽
            const maxIdLen = Math.max(...items.map(i => i.itemId.length));
            items.forEach(it => {
                const padded = it.itemId.padEnd(maxIdLen, ' ');
                console.log(`   ${padded}  × ${it.count}`);
            });
        }
    } else {
        printResult(resp, '查询背包');
    }
}

async function handlePerfTest(args) {
    if (!currentToken) { console.log('❌ 未登录'); return; }
    const n = parseInt(args[0] || '100', 10);
    const concurrency = parseInt(args[1] || '10', 10); // 并发数，默认 10
    if (isNaN(n) || n <= 0) { console.log('❌ 次数无效'); return; }
    if (isNaN(concurrency) || concurrency <= 0) { console.log('❌ 并发数无效'); return; }

    console.log(`🏃 压测 /api/me × ${n} 次 (并发 ${concurrency})...`);

    const start = Date.now();
    let success = 0, fail = 0;
    let http401 = 0, http404 = 0, httpOther = 0, netErr = 0;

    // 用 worker pool 控制并发
    async function worker() {
        while (true) {
            // 用闭包取任务序号
            let idx;
            {
                idx = counter++;
            }
            if (idx >= n) break;
            try {
                const resp = await get('/api/me');
                if (resp.body.code === 0) {
                    success++;
                } else {
                    fail++;
                    if (resp.status === 401) http401++;
                    else if (resp.status === 404) http404++;
                    else httpOther++;
                }
            } catch (e) {
                fail++;
                netErr++;
            }
        }
    }

    let counter = 0;
    const workers = [];
    for (let i = 0; i < Math.min(concurrency, n); i++) {
        workers.push(worker());
    }
    await Promise.all(workers);

    const elapsed = Date.now() - start;
    const qps = elapsed > 0 ? Math.round((n / elapsed) * 1000) : '∞';
    const avg = elapsed > 0 ? (elapsed / n).toFixed(2) : '0.00';

    console.log(`
📊 压测结果:
  总请求:    ${n}
  并发:      ${Math.min(concurrency, n)}
  ✅ 成功:    ${success}
  ❌ 失败:    ${fail}
     ├─ 401 (Token无效): ${http401}
     ├─ 404:             ${http404}
     ├─ 其它HTTP:        ${httpOther}
     └─ 网络错误:        ${netErr}
  ⏱ 总耗时:  ${elapsed}ms
  📈 平均:    ${avg}ms/次
  🚀 QPS:     ${qps}
    `);
}

function printHelp() {
    console.log(`
📋 可用命令:
  ─── 账号 ───
  register <u> <p>         注册新用户
  login <u> <p>            登录（自动缓存 Token）
  logout                   登出（Redis 销毁 Token）
  me                       查看当前登录用户（/api/me 带 Token）
  user <uid>               查询指定用户信息（/api/user/{uid}）
  whoami                   显示本地缓存的 Token/UID

  ─── 测试 ───
  health                   健康检查（GET /）
  perf [次数]              压测 /api/me（默认 100 次）

  ─── 其它 ───
  quit / exit              退出
  help                     显示帮助

当前 Token 状态: ${currentToken ? '✅ 已登录 (' + currentUsername + ')' : '❌ 未登录'}
`);
}

// ═══════════════════════════════════════════════
//  REPL 主循环
// ═══════════════════════════════════════════════
const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

console.log('╔══════════════════════════════════════════╗');
console.log('║  mssgServer HTTP API 测试客户端           ║');
console.log('║  短连接模式 · 零依赖                      ║');
console.log('╚══════════════════════════════════════════╝');
console.log(`  目标: ${opts.baseUrl}`);
console.log(`  时间: ${new Date().toLocaleString()}\n`);

async function mainLoop() {
    // 自动登录
    if (opts.username && opts.password) {
        console.log(`🔄 自动登录: ${opts.username}`);
        await handleLogin([opts.username, opts.password]);
    }

    rl.setPrompt('api> ');
    console.log('输入 help 查看命令列表\n');
    rl.prompt();

    rl.on('line', async (line) => {
        const input = line.trim();
        if (!input) { rl.prompt(); return; }

        const [cmd, ...args] = input.split(/\s+/);

        try {
            switch (cmd.toLowerCase()) {
                case 'quit':
                case 'exit':
                    console.log('👋 bye~');
                    rl.close();
                    return;
                case 'help':
                case '?':
                    printHelp();
                    break;
                case 'register':  await handleRegister(args); break;
                case 'login':     await handleLogin(args); break;
                case 'logout':    await handleLogout(); break;
                case 'me':        await handleMe(); break;
                case 'user':      await handleUserInfo(args); break;
                case 'whoami':    await handleWhoami(); break;
                case 'health':    await handleHealth(); break;
                case 'perf':      await handlePerfTest(args); break;

                // ─── 道具模块 ───
                case 'additem':
                case 'add':       await handleItemAdd(args); break;
                case 'consume':
                case 'useitem':   await handleItemConsume(args); break;
                case 'inventory':
                case 'bag':
                case 'items':     await handleItemGet(); break;

                default:
                    console.log('❌ 未知命令，输入 help 查看帮助');
            }
        } catch (e) {
            console.log('❌ 请求异常:', e.message);
        }
        rl.prompt();
    });
}

mainLoop();