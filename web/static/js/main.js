// --- gRPC official client imports ---
import { HelloRequest } from "./pb/hello_pb.js";
import { HelloServiceClient } from "./pb/hello_grpc_web_pb.js";

import { PingRequest } from "./pb/ping_pb.js";
import { PingServiceClient } from "./pb/ping_grpc_web_pb.js";

// --- google Empty ---
import { Empty } from "google-protobuf/google/protobuf/empty_pb.js";

// --- ENV ---
const PORT = process.env.PORT || "1111";
const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || "2222";
const PROTOCOL = window.location.protocol;
const HOSTNAME = window.location.hostname;
const GRPC_HOST = `${PROTOCOL}//${HOSTNAME}:${GRPC_WEB_PORT}`;

// --- gRPC clients ---
const helloClient = new HelloServiceClient(GRPC_HOST, null, null);
const pingClient = new PingServiceClient(GRPC_HOST, null, null);

/**
 * SayHello
 */
export function sayHello(name, token) {
    return new Promise((resolve, reject) => {
        console.log("Say hello.")
        const req = new HelloRequest();
        req.setName(name);

        helloClient.sayHello(
            req,
            { "Authorization": "Bearer " + token}, (err, resp) => {
            if (err) {
                reject(err);
            } else {
                resolve(resp);
            }
        });
    });
}

/**
 * Ping
 */
export function ping() {
    return new Promise((resolve, reject) => {
        console.log("Ping.")
        const req = new Empty();

        pingClient.ping(req, {}, (err, resp) => {
            if (err) {
                reject(err);
            } else {
                resolve(resp);
            }
        });
    });
}

// --------------------------
// WebSocket chat
// --------------------------
export function initChat(cfg) {
    const myUserId = cfg.userId;
    const token = cfg.token;

    if (!token) {
        cfg.status.textContent = "❌ Token not found. Please log in.";
        return;
    }

    const host = `ws://${HOSTNAME}:${PORT}/ws`;
    const url = host + "?token=" + token;

    const socket = new WebSocket(url);

    socket.onopen = () =>
        (cfg.status.textContent = "🟢 Connected (" + host + ")");
    socket.onclose = () => (cfg.status.textContent = "🔴 Disconnected");
    socket.onerror = () => (cfg.status.textContent = "❌ Error");

    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);

            if (msg.type === "message") {
                const isMine =
                    String(msg.userId) === String(myUserId);
                const text = isMine
                    ? msg.text
                    : `${msg.userId}: ${msg.text}`;
                addMessage(text, isMine);
            }
        } catch (e) {
            console.error("Bad JSON:", e);
        }
    };

    cfg.sendBtn.onclick = sendMessage;
    cfg.input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            sendMessage();
        }
    });

    function sendMessage() {
        const text = cfg.input.value.trim();
        if (!text || socket.readyState !== WebSocket.OPEN) return;

        socket.send(
            JSON.stringify({
                type: "message",
                text,
            })
        );

        cfg.input.value = "";
    }

    function addMessage(text, isMine) {
        const div = document.createElement("div");
        div.textContent = text;

        if (isMine) {
            div.classList.add("chat-my-message");
        }

        cfg.messages.appendChild(div);
        cfg.messages.scrollTop = cfg.messages.scrollHeight;
    }
}
