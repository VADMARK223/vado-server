const GRPC_WEB_URL = 'http://localhost:8090';

// Простая проверка health endpoint
async function testHealth() {
    const result = document.getElementById('result');
    result.textContent = 'Loading...\n';

    try {
        const response = await fetch(GRPC_WEB_URL + '/health');
        const text = await response.text();
        result.textContent = `✅ Health Check!!!\nStatus: ${response.status}\nResponse: ${text}`;
    } catch (e) {
        result.textContent = '❌ Error: ' + e.message;
    }
}

// Вызов gRPC метода Ping
async function testPing() {
    const result = document.getElementById('result');
    result.textContent = 'Calling Ping...\n';

    try {
        // Empty message для Ping (размер 0)
        const body = new Uint8Array([0, 0, 0, 0, 0]);

        const response = await fetch(GRPC_WEB_URL + '/ping.PingService/Ping', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/grpc-web+proto',
                'X-Grpc-Web': '1',
                'X-User-Agent': 'grpc-web-javascript/0.1'
            },
            body: body
        });

        result.textContent = `✅ Ping Response\n`;
        result.textContent += `Status: ${response.status}\n`;
        result.textContent += `Content-Type: ${response.headers.get('content-type')}\n\n`;

        // Читаем бинарный ответ
        const arrayBuffer = await response.arrayBuffer();
        const bytes = new Uint8Array(arrayBuffer);

        result.textContent += `Response bytes (${bytes.length}): ${Array.from(bytes).join(', ')}\n\n`;

        // Пытаемся распарсить как text (для отладки)
        const text = new TextDecoder().decode(bytes);
        result.textContent += `As text: ${text}`;

        if (response.ok) {
            result.textContent += '\n\n✅ Server responded successfully!';
        }
    } catch (e) {
        result.textContent = '❌ Error: ' + e.message + '\n' + e.stack;
    }
}