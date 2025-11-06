import { grpc } from "@improbable-eng/grpc-web";
import { HelloRequest, HelloResponse } from "./pb/hello_pb";
import { HelloService } from "./pb/hello_pb_service";

import { PingResponse } from "./pb/ping_pb";
import { PingService } from "./pb/ping_pb_service";

import { Empty } from "google-protobuf/google/protobuf/empty_pb";

declare const process: any;

const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || '1111'
const HOST = 'http://localhost:' + GRPC_WEB_PORT;

export function sayHello(name: string): Promise<HelloResponse.AsObject> {
    const req = new HelloRequest();
    req.setName(name);

    const md = new grpc.Metadata();
    const transport = grpc.CrossBrowserHttpTransport({ withCredentials: true });

    return new Promise((resolve, reject) => {
        grpc.unary(HelloService.SayHello, {
            request: req,
            host: HOST,
            metadata: md,
            transport: transport,
            onEnd: (res) => {
                if (res.status === grpc.Code.OK && res.message) {
                    resolve(res);
                } else {
                    reject(new Error(res.statusMessage || "gRPC error " + res.status));
                }
            },
        });
    });
}

export function ping(): Promise<PingResponse.AsObject> {
    const req = new Empty();

    return new Promise((resolve, reject) => {
        grpc.unary(PingService.Ping, {
            request: req,
            host: HOST,
            onEnd: (res) => {
                if (res.status === grpc.Code.OK && res.message) {
                    resolve(res);
                } else {
                    reject(new Error(res.statusMessage || "gRPC error " + res.status));
                }
            },
        });
    });
}