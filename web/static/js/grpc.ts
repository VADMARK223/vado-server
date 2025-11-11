import { grpc } from "@improbable-eng/grpc-web";
import { HelloRequest, HelloResponse } from "./pb/hello_pb";
import { HelloService } from "./pb/hello_pb_service";

import { PingResponse } from "./pb/ping_pb";
import { PingService } from "./pb/ping_pb_service";

import { Empty } from "google-protobuf/google/protobuf/empty_pb";

declare const process: any;

const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || '1111'
// const HOST = 'http://localhost:' + GRPC_WEB_PORT;
const HOST = `${window.location.protocol}//${window.location.hostname}:${GRPC_WEB_PORT}`;
const defaultTransport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

export function sayHello(name: string): Promise<HelloResponse> {
    const req = new HelloRequest();
    req.setName(name);

    const md = new grpc.Metadata();

    console.log("HOST:" + HOST)

    return new Promise((resolve, reject) => {
        grpc.unary(HelloService.SayHello, {
            request: req,
            host: HOST,
            metadata: md,
            transport: defaultTransport,
            onEnd: (res) => {
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
            onEnd: (res) => {
                if (res.status === grpc.Code.OK && res.message) {
                    resolve(res.message as PingResponse);
                } else {
                    reject(new Error(res.statusMessage || "gRPC error " + res.status));
                }
            },
        });
    });
}