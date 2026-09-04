const ws = new WebSocket("ws://127.0.0.1:8003/");

ws.onopen = () => {
    console.log("连接成功");
    ws.send('{"Name":"account.login"}');
};

ws.onmessage = (event) => {
    console.log("收到消息:", event.data);
};

ws.onclose = () => {
    console.log("连接已关闭");
};

ws.onerror = (error) => {
    console.error("发生错误:", error);
};

