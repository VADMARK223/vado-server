package chat

import (
	"context"
	"fmt"
	"sync"
	"time"
	"vado_server/api/pb/chat"
)

var clientColor = []string{"#FF5733", "#33FF57", "#3357FF", "#FF33A1", "#33FFF5"}
var clientIndex = 0

type Server struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[chat.ChatService_ChatStreamServer]*Client
}

func NewChatService() *Server {
	return &Server{clients: make(map[chat.ChatService_ChatStreamServer]*Client)}
}

func (s *Server) SendMessage(_ context.Context, msg *chat.ChatMessage) (*chat.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		var messageType chat.MessageType
		if s.clients[client].id == msg.Id {
			messageType = chat.MessageType_MESSAGE_SELF
		} else {
			messageType = chat.MessageType_MESSAGE_USER
		}
		messageWithTime(msg, messageType)
		err := client.Send(msg)
		if err != nil {
			delete(s.clients, client)
		}
	}
	return &chat.Empty{}, nil
}

func (s *Server) ChatStream(req *chat.ChatStreamRequest, stream chat.ChatService_ChatStreamServer) error {
	s.mu.Lock()
	color := clientColor[clientIndex%len(clientColor)]
	clientIndex++

	s.clients[stream] = &Client{
		stream: stream,
		id:     req.Id,
		color:  color,
	}
	s.broadcastSystemMessage("Новый участник вошел", len(s.clients))
	s.mu.Unlock()

	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.clients, stream)
	s.broadcastSystemMessage("Участник покинул", len(s.clients))
	s.mu.Unlock()

	return nil
}

func (s *Server) broadcastSystemMessage(text string, users int) {
	msg := &chat.ChatMessage{
		User:  "System",
		Text:  fmt.Sprintf("%s | Сейчас в чате: %d", text, users),
		Color: "#888888",
	}

	for _, c := range s.clients {
		messageWithTime(msg, chat.MessageType_MESSAGE_USER)
		errSend := c.stream.Send(msg)
		if errSend != nil {
			fmt.Println("Error send message:" + errSend.Error())
		}
	}

}

func messageWithTime(msg *chat.ChatMessage, messageType chat.MessageType) {
	msg.Timestamp = time.Now().Unix()
	msg.Type = messageType
}
