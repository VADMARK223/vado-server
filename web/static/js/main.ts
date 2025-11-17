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
    status: HTMLElement,
    messages: HTMLElement,
    input: HTMLInputElement,
    sendBtn: HTMLButtonElement
}) {
    console.log("Init chat, host: ");
    const socket = new WebSocket("ws://localhost:5555/ws");

    socket.onopen = () => cfg.status.textContent = "🟢 Connected";
    socket.onclose = () => cfg.status.textContent = "🔴 Disconnected";
    socket.onerror = () => cfg.status.textContent = "❌ Error";

    socket.onmessage = (e) => addMessage(e.data);

    cfg.sendBtn.onclick = sendMessage;
    cfg.input.onkeypress = (e) => e.key === "Enter" && sendMessage();

    function sendMessage() {
        const text = cfg.input.value.trim();
        if (!text) return;

        socket.send(text);
        cfg.input.value = "";
    }

    function addMessage(text: string) {
        const div = document.createElement("div");
        div.textContent = text;
        cfg.messages.appendChild(div);
        cfg.messages.scrollTop = cfg.messages.scrollHeight;
    }
}