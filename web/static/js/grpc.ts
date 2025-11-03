import { PingServiceClient } from './pb/PingServiceClientPb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

declare const process: any;

const GRPC_WEB_PORT = process.env.GRPC_WEB_PORT || '1111'
const GRPC_WEB_URL = 'http://localhost:' + GRPC_WEB_PORT;
const pingClient = new PingServiceClient(GRPC_WEB_URL, null, null);

export async function pingServer(): Promise<Empty> {
    return new Promise((resolve, reject) => {
        pingClient.ping(new Empty(), {}, (err, resp) => {
            if (err) reject(err);
            else resolve(resp);
        });
    });
}