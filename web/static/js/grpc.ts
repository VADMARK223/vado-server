import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

import { PingServiceClient } from './pb/PingServiceClientPb';
import { PingResponse } from './pb/ping_pb'

import { HelloServiceClient } from './pb/HelloServiceClientPb';
import { HelloRequest, HelloResponse } from './pb/hello_pb';

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

export async function sayHello(name: string): Promise<HelloResponse> {
    const req = new HelloRequest();
    req.setName(name);

    return new Promise<HelloResponse>((resolve, reject) => {
        helloClient.sayHello(req, {}, (err, resp) => {
            if (err || !resp) reject(err);
            else resolve(resp);
        });
    });
}