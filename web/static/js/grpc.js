import {PingServiceClient} from './pb/ping_grpc_web_pb.js';
import {Empty} from 'google-protobuf/google/protobuf/empty_pb.js';

const GRPC_WEB_URL = 'http://localhost:8090';
const pingClient = new PingServiceClient(GRPC_WEB_URL, null, null);

export async function pingServer() {
    return new Promise((resolve, reject) => {
        pingClient.ping(new Empty(), {}, (err, resp) => {
            if (err) reject(err);
            else resolve(resp);
        });
    });
}