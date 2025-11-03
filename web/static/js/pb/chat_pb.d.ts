import * as jspb from 'google-protobuf'



export class ChatMessage extends jspb.Message {
  getId(): string;
  setId(value: string): ChatMessage;

  getUser(): User | undefined;
  setUser(value?: User): ChatMessage;
  hasUser(): boolean;
  clearUser(): ChatMessage;

  getText(): string;
  setText(value: string): ChatMessage;

  getTimestamp(): number;
  setTimestamp(value: number): ChatMessage;

  getType(): MessageType;
  setType(value: MessageType): ChatMessage;

  getUsersCount(): number;
  setUsersCount(value: number): ChatMessage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatMessage.AsObject;
  static toObject(includeInstance: boolean, msg: ChatMessage): ChatMessage.AsObject;
  static serializeBinaryToWriter(message: ChatMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatMessage;
  static deserializeBinaryFromReader(message: ChatMessage, reader: jspb.BinaryReader): ChatMessage;
}

export namespace ChatMessage {
  export type AsObject = {
    id: string,
    user?: User.AsObject,
    text: string,
    timestamp: number,
    type: MessageType,
    usersCount: number,
  }
}

export class User extends jspb.Message {
  getId(): number;
  setId(value: number): User;

  getUsername(): string;
  setUsername(value: string): User;

  getColor(): string;
  setColor(value: string): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: number,
    username: string,
    color: string,
  }
}

export class ChatStreamRequest extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): ChatStreamRequest;
  hasUser(): boolean;
  clearUser(): ChatStreamRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatStreamRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChatStreamRequest): ChatStreamRequest.AsObject;
  static serializeBinaryToWriter(message: ChatStreamRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatStreamRequest;
  static deserializeBinaryFromReader(message: ChatStreamRequest, reader: jspb.BinaryReader): ChatStreamRequest;
}

export namespace ChatStreamRequest {
  export type AsObject = {
    user?: User.AsObject,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

export enum MessageType { 
  MESSAGE_UNKNOWN = 0,
  MESSAGE_USER = 1,
  MESSAGE_SYSTEM = 2,
  MESSAGE_SELF = 3,
}
