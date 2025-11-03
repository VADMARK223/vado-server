import { PingServiceClient } from './pb/PingServiceClientPb';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

const GRPC_WEB_URL = 'http://localhost:8090';
const pingClient = new PingServiceClient(GRPC_WEB_URL, null, null);

export async function pingServer(): Promise<Empty> {
    return new Promise((resolve, reject) => {
        pingClient.ping(new Empty(), {}, (err, resp) => {
            if (err) reject(err);
            else resolve(resp);
        });
    });
}