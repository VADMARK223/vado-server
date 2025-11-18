// @ts-ignore
import { grpc } from "@improbable-eng/grpc-web";
import { HelloRequest, HelloResponse } from "./pb/hello_pb";
import { HelloService } from "./pb/hello_pb_service";

import { PingResponse } from "./pb/ping_pb";
import { PingService } from "./pb/ping_pb_service";

import { Empty } from "google-protobuf/google/protobuf/empty_pb";

declare const process: any;

const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || '1111'
const HOST = `${window.location.protocol}//${window.location.hostname}:${GRPC_WEB_PORT}`;
const defaultTransport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

export function sayHello(name: string): Promise<HelloResponse> {
    const req = new HelloRequest();
    req.setName(name);

    const md = new grpc.Metadata();

    console.log("gRPC host: " + HOST)

    return new Promise((resolve, reject) => {
        grpc.unary(HelloService.SayHello, {
            request: req,
            host: HOST,
            metadata: md,
            transport: defaultTransport,
            onEnd: (res: { status: string; message: HelloResponse; statusMessage: any; }) => {
                if (res.status === grpc.Code.OK && res.message) {
                    resolve(res.message as HelloResponse);
                } else {
                    reject(new Error(res.statusMessage || "gRPC error " + res.status));
                }
            },
        });
    });
}

export function ping(): Promise<PingResponse> {
    const req = new Empty();

    console.log("HOST:" + HOST)

    return new Promise((resolve, reject) => {
        grpc.unary(PingService.Ping, {
            request: req,
            host: HOST,
            onEnd: (res: { status: string; message: PingResponse; statusMessage: any; }) => {
                if (res.status === grpc.Code.OK && res.message) {
                    resolve(res.message as PingResponse);
                } else {
                    reject(new Error(res.statusMessage || "gRPC error " + res.status));
                }
            },
        });
    });
}

export function initChat(cfg: {
    userId: string
    status: HTMLElement,
    messages: HTMLElement,
    input: HTMLInputElement,
    sendBtn: HTMLButtonElement
}) {
    const myUserId = cfg.userId;
    console.log("BUNDLE: my user id: " + myUserId)
    const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlcyI6WyJ1c2VyIl0sImlzcyI6InZhZG8tc2VydmVyIiwic3ViIjoiYWNjZXNzIiwiZXhwIjoxNzY0MDI3OTI0LCJuYmYiOjE3NjM0MjMxMjQsImlhdCI6MTc2MzQyMzEyNH0.sEaDHzD9UzYyXPk3Qsi0Wwlc9HEeomhVu12j98AHauI"
    const host = "ws://localhost:5555/ws"
    const url = host + "?token=" + token
    const socket = new WebSocket(url);

    socket.onopen = () => cfg.status.textContent = "🟢 Connected (" + host + ")";
    socket.onclose = () => cfg.status.textContent = "🔴 Disconnected";
    socket.onerror = () => cfg.status.textContent = "❌ Error";

    // socket.onmessage = (e) => addMessage(e.data);
    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            if (msg.type === "message") {
                const isMine = String(msg.userId) === String(myUserId)
                addMessage(`${msg.userId}: ${msg.text}`, isMine);
            } else {
                console.log("Other msg:", msg);
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

        const packet = {
            type: "message",
            text: text,
        };

        socket.send(JSON.stringify(packet));
        cfg.input.value = "";
    }

    function addMessage(text:string, isMine = false) {
        const div = document.createElement("div");
        div.textContent = text;

        if (isMine) {
            div.classList.add("chat-my-message")
        }

        cfg.messages.appendChild(div);
        cfg.messages.scrollTop = cfg.messages.scrollHeight;
    }
}