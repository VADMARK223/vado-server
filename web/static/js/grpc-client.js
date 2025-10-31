
// Инициализируем proto объект
window.proto = window.proto || {};
window.protoModules = window.protoModules || {};

// Простая имитация require для браузера
window.require = function(name) {
    if (name === 'google-protobuf') {
        return window.jspb;
    }
    if (name === 'grpc-web') {
        return window.grpc.web;
    }

    if (name.startsWith('./') || name.startsWith('../')) {
        const moduleName = name.replace('./', '').replace('../', '');
        // Возвращаем уже загруженный модуль
        return window.protoModules[moduleName] || window.proto;
    }

    // Для proto файлов возвращаем window.proto
    return window.proto;
};

// Перехватываем exports для сбора всех proto классов
let currentModule = null;

Object.defineProperty(window, 'exports', {
    get: function() {
        if (!currentModule) {
            currentModule = {};
        }
        return currentModule;
    },
    set: function(value) {
        if (value && typeof value === 'object') {
            Object.assign(window.proto, value);
            if (currentModule) {
                Object.assign(currentModule, value);
            }
        }
    },
    configurable: true
});

Object.defineProperty(window, 'module', {
    get: function() {
        return {
            get exports() {
                if (!currentModule) {
                    currentModule = {};
                }
                return currentModule;
            },
            set exports(value) {
                if (value && typeof value === 'object') {
                    Object.assign(window.proto, value);
                    currentModule = value;
                }
            }
        };
    },
    set: function(value) {
        // Игнорируем
    },
    configurable: true
});

// Функция для регистрации загруженного модуля
window.registerProtoModule = function(name, exports) {
    window.protoModules[name] = exports;
    Object.assign(window.proto, exports);
};

// Базовая реализация gRPC-Web
window.grpc = window.grpc || {};
window.grpc.web = window.grpc.web || {};

// GrpcWebClientBase
window.grpc.web.GrpcWebClientBase = class GrpcWebClientBase {
    constructor(options) {
        this.hostname_ = options.hostname || '';
        this.credentials_ = options.credentials;
        this.options_ = options.options;
    }

    rpcCall(method, request, metadata, methodDescriptor, callback) {
        const hostname = this.hostname_;
        const methodPath = methodDescriptor.name;

        // Сериализуем request
        const serialized = request.serializeBinary();

        // Создаём gRPC-Web frame (5 байт header + data)
        const frame = new Uint8Array(5 + serialized.length);
        frame[0] = 0; // compressed flag
        frame[1] = (serialized.length >> 24) & 0xFF;
        frame[2] = (serialized.length >> 16) & 0xFF;
        frame[3] = (serialized.length >> 8) & 0xFF;
        frame[4] = serialized.length & 0xFF;
        frame.set(serialized, 5);

        fetch(hostname + '/' + methodPath, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/grpc-web+proto',
                'X-Grpc-Web': '1',
                'X-User-Agent': 'grpc-web-javascript/0.1',
                ...metadata
            },
            body: frame
        })
            .then(response => response.arrayBuffer())
            .then(arrayBuffer => {
                const bytes = new Uint8Array(arrayBuffer);

                // Парсим gRPC-Web response (пропускаем 5 байт header)
                if (bytes.length < 5) {
                    throw new Error('Invalid response');
                }

                const messageLength = (bytes[1] << 24) | (bytes[2] << 16) | (bytes[3] << 8) | bytes[4];
                const messageBytes = bytes.slice(5, 5 + messageLength);

                // Десериализуем response
                const ResponseType = methodDescriptor.responseType;
                const response = ResponseType.deserializeBinary(messageBytes);

                callback(null, response);
            })
            .catch(error => {
                callback({
                    code: 2, // UNKNOWN
                    message: error.message
                });
            });
    }
};

window.grpc.web.MethodDescriptor = class MethodDescriptor {
    constructor(name, methodType, requestType, responseType, requestSerializeFn, responseDeserializeFn) {
        this.name = name;
        this.methodType = methodType;
        this.requestType = requestType;
        this.responseType = responseType;
        this.requestSerializeFn = requestSerializeFn;
        this.responseDeserializeFn = responseDeserializeFn;
    }
};

window.grpc.web.MethodType = {
    UNARY: 'unary',
    SERVER_STREAMING: 'server_streaming',
    BIDI_STREAMING: 'bidi_streaming'
};

window.grpc.web.AbstractClientBase = window.grpc.web.GrpcWebClientBase;