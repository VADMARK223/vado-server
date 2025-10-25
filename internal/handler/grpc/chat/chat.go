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
	mu sync.Mutex
	//clients map[chat.ChatService_ChatStreamServer]*Client
	clients map[uint64]*Client
}

func NewChatService() *Server {
	return &Server{clients: make(map[uint64]*Client)}
}

func (s *Server) ChatStream(req *chat.ChatStreamRequest, stream chat.ChatService_ChatStreamServer) error {
	s.mu.Lock()
	color := clientColor[clientIndex%len(clientColor)]
	clientIndex++

	userID := req.User.Id

	s.clients[userID] = &Client{
		stream: stream,
		user:   &chat.User{Id: userID, Username: req.User.Username, Color: color},
	}
	s.broadcastSystemMessage(req.User.Id, fmt.Sprintf("Новый участник: %s", req.User.Username), len(s.clients))
	s.mu.Unlock()

	<-stream.Context().Done()

	s.mu.Lock()
	delete(s.clients, userID)
	s.broadcastSystemMessage(req.User.Id, fmt.Sprintf("Участник покинул: %s", req.User.Username), len(s.clients))
	s.mu.Unlock()

	return nil
}

func (s *Server) SendMessage(_ context.Context, msg *chat.ChatMessage) (*chat.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sender, ok := s.clients[msg.User.Id]
	if !ok {
		return nil, fmt.Errorf("unknown sender with id %d", msg.User.Id)
	}

	color := sender.user.Color

	for id, client := range s.clients {
		var messageType chat.MessageType
		if id == msg.User.Id {
			messageType = chat.MessageType_MESSAGE_SELF
		} else {
			messageType = chat.MessageType_MESSAGE_USER
		}
		messageWithTime(msg, messageType)
		msg.User.Color = color
		err := client.stream.Send(msg)
		if err != nil {
			delete(s.clients, id)
		}
	}
	return &chat.Empty{}, nil
}

func (s *Server) broadcastSystemMessage(userId uint64, text string, usersCount int) {
	msg := &chat.ChatMessage{
		User: &chat.User{Id: userId, Username: "System", Color: "#888888"},
		Text: fmt.Sprintf("%s", text),
	}

	for _, c := range s.clients {
		messageWithTime(msg, chat.MessageType_MESSAGE_SYSTEM)
		msg.UsersCount = uint32(usersCount)
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
