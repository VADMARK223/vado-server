package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"
	"vado_server/api/pb/chat"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	stream chat.ChatService_ChatStreamServer
	user   *chat.User
}

var clientColor = []string{"#FF5733", "#33FF57", "#3357FF", "#FF33A1", "#33FFF5"}
var clientIndex = 0

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[uint64]*Client
}

func New() *ChatServer {
	return &ChatServer{clients: make(map[uint64]*Client)}
}

func (s *ChatServer) ChatStream(req *chat.ChatStreamRequest, stream chat.ChatService_ChatStreamServer) error {
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

func (s *ChatServer) SendMessage(_ context.Context, msg *chat.ChatMessage) (*chat.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sender, ok := s.clients[msg.User.Id]
	if !ok {
		return nil, fmt.Errorf("unknown sender with id %d", msg.User.Id)
	}

	senderUser := sender.user

	for id, client := range s.clients {
		our := cloneMsgFor(senderUser, msg, id == senderUser.Id)

		if err := client.stream.Send(our); err != nil {
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.Canceled:
				fmt.Println("Client canceled the stream.")
				delete(s.clients, id)
			case codes.Unavailable:
				fmt.Println("Client unavailable.")
				delete(s.clients, id)
			default:
				fmt.Println("Client error:", err)
				delete(s.clients, id)
			}
		}
	}
	return &chat.Empty{}, nil
}

func cloneMsgFor(sender *chat.User, in *chat.ChatMessage, isSelf bool) *chat.ChatMessage {
	mt := chat.MessageType_MESSAGE_USER

	if isSelf {
		mt = chat.MessageType_MESSAGE_SELF
	}

	return &chat.ChatMessage{
		User: &chat.User{
			Id:       sender.Id,
			Username: sender.Username,
			Color:    sender.Color,
		},
		Text:      in.Text,
		Timestamp: time.Now().Unix(),
		Type:      mt,
	}
}

func (s *ChatServer) broadcastSystemMessage(userId uint64, text string, usersCount int) {
	for id, c := range s.clients {
		out := &chat.ChatMessage{
			User:       &chat.User{Id: userId, Username: "System", Color: "#888888"},
			Text:       text,
			Timestamp:  time.Now().Unix(),
			Type:       chat.MessageType_MESSAGE_SYSTEM,
			UsersCount: uint32(usersCount),
		}
		if err := c.stream.Send(out); err != nil {
			st, _ := status.FromError(err)
			if st.Code() == codes.Canceled || st.Code() == codes.Unavailable {
				delete(s.clients, id)
				continue
			}
			delete(s.clients, id)
		}
	}
}
