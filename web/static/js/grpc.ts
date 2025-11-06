import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

import { PingServiceClient } from './pb/PingServiceClientPb';
import { PingResponse } from './pb/ping_pb'

import { HelloServiceClient } from './pb/HelloServiceClientPb';
import { HelloRequest, HelloResponse } from './pb/hello_pb';
import { grpc } from "@improbable-eng/grpc-web";

declare const process: any;

const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || '1111'
const GRPC_WEB_URL = 'http://localhost:' + GRPC_WEB_PORT;

const pingClient = new PingServiceClient(GRPC_WEB_URL, null, null);
const helloClient = new HelloServiceClient(GRPC_WEB_URL, null, null);

export async function pingServer(): Promise<PingResponse> {
    return new Promise<PingResponse>((resolve, reject) => {
        pingClient.ping(new Empty(), {}, (err, resp) => {
            if (err || !resp) reject(err);
            else resolve(resp);
        });
    });
}


export async function sayHello(name: string, token: string): Promise<HelloResponse.AsObject> {
    const req = new HelloRequest();
    req.setName(name);

    const metadata = new grpc.Metadata();
    metadata.set("authorization", `Bearer ${token}`);

    return new Promise((resolve, reject) => {
        grpc.unary(HelloServiceClient.SayHello, {
            request: req,
            host: GRPC_WEB_URL,
            metadata,
            onEnd: (res) => {
                const { status, message, statusMessage } = res;
                if (status === grpc.Code.OK && message) {
                    resolve(message.toObject());
                } else {
                    reject(new Error(statusMessage || "Unknown error"));
                }
            },
        });
    });
}


export async function sayHelloOld(name: string): Promise<HelloResponse> {
    const req = new HelloRequest();
    req.setName(name);

    const metadata = new grpc.Metadata();
    metadata.set("authorization", `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1NCwicm9sZXMiOlsidXNlciJdLCJpc3MiOiJ2YWRvLXBpbmciLCJzdWIiOiJhY2Nlc3MiLCJleHAiOjE3NjI1MDYwNzksImlhdCI6MTc2MjQxOTY3OX0.qj4p8ltCakA8dNUknAEsc0k2fqWsxjbXfjpPU0YDPvc`);

    return new Promise<HelloResponse>((resolve, reject) => {
        helloClient.sayHello(req, metadata, (err, resp) => {
            if (err || !resp) reject(err);
            else resolve(resp);
        });
    });
}