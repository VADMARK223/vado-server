
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

        const arrayBuffer = await response.arrayBuffer();
        const bytes = new Uint8Array(arrayBuffer);

        result.textContent += `Response bytes (${bytes.length}): ${Array.from(bytes).join(', ')}\n\n`;

        const text = new TextDecoder().decode(bytes);
        result.textContent += `As text: ${text}`;

        if (response.ok) {
            result.textContent += '\n\n✅ Server responded successfully!';
        }
    } catch (e) {
        result.textContent = '❌ Error: ' + e.message + '\n' + e.stack;
    }
}

// Вызов gRPC метода SayHello
function testSayHello() {
    const result = document.getElementById('result');
    result.textContent = 'Calling SayHello...\n';

    try {
        // Отладочный вывод
        console.log('proto object:', proto);
        console.log('Available keys:', Object.keys(proto));

        // Проверяем наличие классов
        if (!proto.HelloRequest) {
            result.textContent = '❌ HelloRequest not found in proto object\n';
            result.textContent += 'Available: ' + Object.keys(proto).join(', ');
            return;
        }

        if (!proto.HelloServiceClient) {
            result.textContent = '❌ HelloServiceClient not found in proto object\n';
            result.textContent += 'Available: ' + Object.keys(proto).join(', ');
            return;
        }

        // Создаём клиент
        const client = new proto.HelloServiceClient(GRPC_WEB_URL, null, null);

        // Создаём request
        const request = new proto.HelloRequest();
        request.setName('Мир');

        // Вызываем метод
        client.sayHello(request, {}, (err, response) => {
            if (err) {
                result.textContent = '❌ Error: ' + err.code + ' - ' + err.message;
                console.error(err);
                return;
            }

            // Получаем ответ
            const message = response.getMessage();
            result.textContent = `✅ SayHello Response\n`;
            result.textContent += `Message: ${message}\n`;
        });

    } catch (e) {
        result.textContent = '❌ Error: ' + e.message + '\n' + e.stack;
        console.error(e);
    }
}