package chat

import "vado_server/api/pb/chat"

type Client struct {
	stream chat.ChatService_ChatStreamServer
	color  string
	name   string
}
