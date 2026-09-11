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
const SERVER_PRESETS = {
    login: 'ws://127.0.0.1:8003',
    chat:  'ws://127.0.0.1:8004'
};

// ═══════════════════════════════════════════════
//  解析命令行参数
// ═══════════════════════════════════════════════
const args = process.argv.slice(2);

function parseArgs(argv) {
    const opts = { server: 'login', username: '', password: '' };
    for (let i = 0; i < argv.length; i++) {
        switch (argv[i]) {
            case '--username': opts.username = argv[++i]; break;
            case '--password': opts.password = argv[++i]; break;
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

const opts = parseArgs(args);
const WS_URL = opts.server.startsWith('ws') 
    ? opts.server 
    : SERVER_PRESETS[opts.server];

// ═══════════════════════════════════════════════
//  状态
// ═══════════════════════════════════════════════
let currentUId = '';
let currentUsername = '';
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
        // 登录响应自动缓存 uid
        if (data.name === 'account.login' && data.code === 0 && data.msg) {
            currentUId = data.msg.uid || '';
            currentUsername = data.msg.username || '';
            console.log(`\n🎯 登录成功! UID=${currentUId} Username=${currentUsername}\n`);
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
  private <toUid> <content>         私聊（自动跨服推到目标节点）

  ─── 群组 ───
  join <groupId>                    加入群组
  leave <groupId>                   退出群组
  members <groupId>                 查看群组成员
  group <groupId> <content>         群组聊天

  ─── 在线 ───
  online                            查看本机在线用户

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
                if (args.length < 2) { console.log('用法: private <toUid> <内容>'); break; }
                send(ws, 'chat.private.send', { toUid: args[0], content: args.slice(1).join(' ') });
                break;

            case 'join':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: join <groupId>'); break; }
                send(ws, 'chat.group.join', { groupId: args[0] });
                break;

            case 'leave':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 1) { console.log('用法: leave <groupId>'); break; }
                send(ws, 'chat.group.leave', { groupId: args[0] });
                break;

            case 'members':
                if (args.length < 1) { console.log('用法: members <groupId>'); break; }
                send(ws, 'chat.group.members', { groupId: args[0] });
                break;

            case 'group':
                if (!currentUId) { console.log('❌ 请先 login'); break; }
                if (args.length < 2) { console.log('用法: group <groupId> <内容>'); break; }
                send(ws, 'chat.group.send', { groupId: args[0], content: args.slice(1).join(' ') });
                break;

            case 'online':
                send(ws, 'chat.online.get', {});
                break;

            default:
                try {
                    const jsonPart = input.slice(cmd.length).trim();
                    const msg = jsonPart ? JSON.parse(jsonPart) : {};
                    send(ws, cmd, msg);
                } catch (e) {
                    console.log('❌ 未知命令，输入 help 查看帮助');
                }
        }
        rl.prompt();
    });
}