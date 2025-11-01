async function deleteUser(id) {
    if (!confirm('Удалить пользователя?')) return;

    const res = await fetch(`/users/${id}`, { method: 'DELETE' });
    if (res.ok) {
        document.getElementById(`user-${id}`).remove();
    } else {
        alert('Ошибка при удалении пользователя');
    }
}

async function deleteTask(id) {
    if (!confirm('Удалить задачу?')) return;

    const res = await fetch(`/tasks/${id}`, { method: 'DELETE' });
    if (res.ok) {
        document.getElementById(`task-${id}`).remove();
    } else {
        alert('Ошибка при удалении задачи');
    }
}

import {PingServiceClient} from './pb/ping_grpc_web_pb.js';
import {Empty} from 'google-protobuf/google/protobuf/empty_pb.js';

// --- Конфигурация клиента ---
const GRPC_WEB_URL = 'http://localhost:8090';
const pingClient = new PingServiceClient(GRPC_WEB_URL, null, null);

// --- UI хелперы ---
function setStatus(message, isError = false) {
    const pre = document.getElementById('result');
    pre.textContent = (isError ? '❌ ' : '✅ ') + message + '\n';
}

// --- API вызовы ---
async function pingServer() {
    return new Promise((resolve, reject) => {
        pingClient.ping(new Empty(), {}, (err, resp) => {
            if (err) reject(err);
            else resolve(resp);
        });
    });
}

// --- Обработчик кнопки ---
export async function onPingClick() {
    setStatus('Выполняется Ping...');
    try {
        const resp = await pingServer();
        setStatus('Сервер ответил: run=' + resp.getRun());
    } catch (err) {
        setStatus(err.message, true);
    }
}

// --- Автоинициализация ---
document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('pingButton').addEventListener('click', onPingClick);
});