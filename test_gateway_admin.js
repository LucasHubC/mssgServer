/**
 * 网关管理接口测试脚本
 * 
 * 用法:
 *   node scripts/test_gateway_admin.js              # 跑全部
 *   node scripts/test_gateway_admin.js broadcast     # 只跑全服广播
 *   node scripts/test_gateway_admin.js activity      # 只跑活动推送
 *   node scripts/test_gateway_admin.js push          # 只跑跨服私信
 *   node scripts/test_gateway_admin.js kick          # 只跑跨服踢人
 *   node scripts/test_gateway_admin.js status        # 只查看节点状态
 *   node scripts/test_gateway_admin.js broadcast push # 组合跑
 */

const http = require('http');

const GATEWAY = { host: '127.0.0.1', port: 9000 };

// ⚠️ 把这里改成你 testWS.js 登录的 uid（push 和 kick 测试会用到）
const TARGET_UID = 'test_a9c2f7853c7c';

// ────────────────────────────────────────
//  通用 HTTP POST / GET
// ────────────────────────────────────────
function post(path, payload) {
    return new Promise((resolve, reject) => {
        const body = JSON.stringify(payload);
        const req = http.request({
            hostname: GATEWAY.host,
            port: GATEWAY.port,
            path,
            method: 'POST',
            headers: {
                'Content-Type': 'application/json; charset=utf-8',
                'Content-Length': Buffer.byteLength(body, 'utf-8')
            }
        }, res => {
            let data = '';
            res.setEncoding('utf-8');
            res.on('data', chunk => data += chunk);
            res.on('end', () => {
                console.log(`  ✅ POST ${path} → ${res.statusCode}`);
                console.log(`     请求: ${JSON.stringify(payload)}`);
                console.log(`     响应: ${data}`);
                resolve();
            });
        });
        req.on('error', reject);
        req.write(body, 'utf-8');
        req.end();
    });
}

function get(path) {
    return new Promise((resolve, reject) => {
        const req = http.request({
            hostname: GATEWAY.host,
            port: GATEWAY.port,
            path,
            method: 'GET'
        }, res => {
            let data = '';
            res.on('data', chunk => data += chunk);
            res.on('end', () => {
                const parsed = JSON.parse(data);
                console.log(`  ✅ GET ${path} → ${res.statusCode}`);
                console.log(`     ${JSON.stringify(parsed, null, 4)}`);
                resolve();
            });
        });
        req.on('error', reject);
        req.end();
    });
}

// ────────────────────────────────────────
//  测试用例
// ────────────────────────────────────────

async function testBroadcast() {
    console.log('\n━━━ ① /admin/broadcast 全服系统公告 ━━━');
    console.log('    所有 loginServer 节点的所有在线用户都会收到这条消息\n');
    await post('/admin/broadcast', {
        name: 'system.notice',
        data: {
            title: '🎉 新版本上线啦',
            content: '今晚 22:00 停服维护 30 分钟，感谢大家支持！',
            time: '2026-09-11 21:00:00'
        }
    });
}

async function testActivity() {
    console.log('\n━━━ ② /admin/broadcast 活动推送 ━━━');
    console.log('    用 broadcast 接口推送活动消息\n');
    await post('/admin/broadcast', {
        name: 'activity.lucky_draw',
        data: {
            activityId: 'act_001',
            prize: 'iPhone 17 Pro',
            endTime: 1750000000,
            rules: '登录即可抽奖，每日 3 次机会'
        }
    });
}

async function testPush() {
    console.log('\n━━━ ③ /admin/push 跨服私信 ━━━');
    console.log(`    目标 UID: ${TARGET_UID}`);
    console.log('    如果该 uid 在线，对应节点会推送给他\n');
    await post('/admin/push', {
        uid: TARGET_UID,
        name: 'chat.private',
        data: {
            from: 'system',
            msg: `你好 ${TARGET_UID}！这是一条跨服私信测试消息 🎁`,
            time: Date.now()
        }
    });
}

async function testKick() {
    console.log('\n━━━ ④ /admin/kick 跨服踢人（慎用！）━━━');
    console.log(`    目标 UID: ${TARGET_UID}`);
    console.log('    会把该 uid 的 WebSocket 连接关掉\n');
    await post('/admin/kick', { uid: TARGET_UID });
}

async function testStatus() {
    console.log('\n━━━ ⑤ /gateway/status 查看节点状态 ━━━\n');
    await get('/gateway/status');
}

// ────────────────────────────────────────
//  路由表
// ────────────────────────────────────────
const TESTS = {
    broadcast: testBroadcast,
    activity:  testActivity,
    push:      testPush,
    kick:      testKick,
    status:    testStatus
};

// ────────────────────────────────────────
//  主入口
// ────────────────────────────────────────
async function main() {
    const args = process.argv.slice(2); // node test.js a b c → ['a', 'b', 'c']

    console.log('============================================');
    console.log('  网关管理接口测试脚本');
    console.log(`  Gateway → ${GATEWAY.host}:${GATEWAY.port}`);
    console.log(`  TARGET  → ${TARGET_UID}`);
    console.log('============================================');

    // ── 没传参数 → 跑全部（kick 默认跳过，踢人太危险）──
    let selected = args;
    if (selected.length === 0) {
        console.log('\n📌 未指定参数，跑全部测试（kick 除外）');
        console.log('   可用参数: broadcast | activity | push | kick | status');
        console.log('   示例: node test_gateway_admin.js broadcast push\n');
        selected = ['status', 'broadcast', 'activity', 'push'];
    }

    // ── 验证参数合法性 ──
    const unknown = selected.filter(k => !TESTS[k]);
    if (unknown.length > 0) {
        console.error(`\n❌ 未知参数: ${unknown.join(', ')}`);
        console.error(`   可用参数: ${Object.keys(TESTS).join(' | ')}`);
        process.exit(1);
    }

    // ── 逐个执行 ──
    try {
        for (const key of selected) {
            await TESTS[key]();
        }
        console.log('\n============================================');
        console.log('  ✅ 全部完成');
        console.log('============================================');
    } catch (e) {
        console.error('\n❌ 执行失败:', e.message);
        console.error('   请确认 Gateway 是否已启动 + Redis 是否已连通');
        process.exit(1);
    }
}

main();