package chat

import (
	"context"
	"sync"
	"vado_server/api/pb/chat"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[chat.ChatService_ChatStreamServer]struct{}
}

func NewChatService() *ChatServer {
	return &ChatServer{clients: make(map[chat.ChatService_ChatStreamServer]struct{})}
}

func (s *ChatServer) SendMessage(_ context.Context, msg *chat.ChatMessage) (*chat.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		err := client.Send(msg)
		if err != nil {
			delete(s.clients, client)
		}
	}
	return &chat.Empty{}, nil
}

func (s *ChatServer) ChatStream(_ *chat.Empty, stream chat.ChatService_ChatStreamServer) error {
	s.mu.Lock()
	s.clients[stream] = struct{}{}
	s.mu.Unlock()

	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.clients, stream)
	s.mu.Unlock()

	return nil
}
