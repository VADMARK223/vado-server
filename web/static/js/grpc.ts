import { grpc } from "@improbable-eng/grpc-web";
import { HelloRequest, HelloResponse } from "./pb/hello_pb";
import { HelloService } from "./pb/hello_pb_service";

declare const process: any;

const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || '1111'
const HOST = 'http://localhost:' + GRPC_WEB_PORT;

export function sayHello(name: string, token: string): Promise<HelloResponse.AsObject> {
    const req = new HelloRequest();
    req.setName(name);

    const md = new grpc.Metadata();
    md.set("authorization", `Bearer ${token}`);

    return new Promise((resolve, reject) => {
        grpc.unary(HelloService.SayHello, {
            request: req,
            host: HOST,
            metadata: md,
            onEnd: (res) => {
                if (res.status === grpc.Code.OK && res.message) {
                    resolve((res.message as any).toObject());
                    resolve(res);
                } else {
                    reject(new Error(res.statusMessage || "gRPC error " + res.status));
                }
            },
        });
    });
}