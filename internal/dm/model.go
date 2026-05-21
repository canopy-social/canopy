package dm

import (
	"context"
	"time"
)

type Conversation struct {
	ID            string     `json:"id"`
	URI           string     `json:"uri"`
	ParticipantA  string     `json:"participant_a"`
	ParticipantB  string     `json:"participant_b"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty"`
	UnreadCountA  int        `json:"unread_count_a"`
	UnreadCountB  int        `json:"unread_count_b"`
	CreatedAt     time.Time  `json:"created_at"`
}

type Message struct {
	ID             string     `json:"id"`
	URI            string     `json:"uri"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	Content        string     `json:"content"`
	IsLocal        bool       `json:"is_local"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type Repository interface {
	CreateConversation(ctx context.Context, conv *Conversation) (*Conversation, error)
	GetConversationByID(ctx context.Context, id string) (*Conversation, error)
	GetConversationByParticipants(ctx context.Context, p1, p2 string) (*Conversation, error)
	ListConversations(ctx context.Context, accountID string, limit, offset int) ([]*Conversation, error)
	UpdateConversation(ctx context.Context, conv *Conversation) error
	CreateMessage(ctx context.Context, msg *Message) (*Message, error)
	GetMessageByID(ctx context.Context, id string) (*Message, error)
	ListMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error)
	MarkMessagesAsRead(ctx context.Context, conversationID, readerID string) error
}
