package chat

import "time"

type MessageLog struct {
	UserID    uint64    `json:"user_id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsSystem  bool      `json:"is_system"`
}
