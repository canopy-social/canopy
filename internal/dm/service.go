package dm

import (
	"context"
	"fmt"
	"time"

	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
)

type AccountLoader interface {
	GetByID(ctx context.Context, id string) (*accounts.Account, error)
}

type Service struct {
	repo          Repository
	accountLoader AccountLoader
}

func NewService(repo Repository, accountLoader AccountLoader) *Service {
	return &Service{
		repo:          repo,
		accountLoader: accountLoader,
	}
}

func (s *Service) SendDM(ctx context.Context, senderID, recipientID, content string) (*Message, error) {
	if senderID == recipientID {
		return nil, fmt.Errorf("cannot send a direct message to yourself")
	}
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}
	_, err := s.accountLoader.GetByID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("sender not found: %w", err)
	}
	_, err = s.accountLoader.GetByID(ctx, recipientID)
	if err != nil {
		return nil, fmt.Errorf("recipient not found: %w", err)
	}
	p1, p2 := senderID, recipientID
	if p1 > p2 {
		p1, p2 = recipientID, senderID
	}
	conv, err := s.repo.GetConversationByParticipants(ctx, p1, p2)
	if err != nil {
		convID := ulid.New()
		convURI := fmt.Sprintf("https://canopy.local/conversations/%s", convID)
		newConv := &Conversation{
			ID:           convID,
			URI:          convURI,
			ParticipantA: p1,
			ParticipantB: p2,
			CreatedAt:    time.Now(),
		}
		conv, err = s.repo.CreateConversation(ctx, newConv)
		if err != nil {
			return nil, fmt.Errorf("creating conversation: %w", err)
		}
	}
	msgID := ulid.New()
	msgURI := fmt.Sprintf("https://canopy.local/messages/%s", msgID)
	now := time.Now()
	msg := &Message{
		ID:             msgID,
		URI:            msgURI,
		ConversationID: conv.ID,
		SenderID:       senderID,
		Content:        content,
		IsLocal:        true,
		CreatedAt:      now,
	}
	createdMsg, err := s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("creating message: %w", err)
	}
	conv.LastMessageAt = &now
	if recipientID == p1 {
		conv.UnreadCountA++
	} else {
		conv.UnreadCountB++
	}
	err = s.repo.UpdateConversation(ctx, conv)
	if err != nil {
		return nil, fmt.Errorf("updating conversation: %w", err)
	}
	return createdMsg, nil
}

func (s *Service) ListConversations(ctx context.Context, accountID string, limit, offset int) ([]*Conversation, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListConversations(ctx, accountID, limit, offset)
}

func (s *Service) ListMessages(ctx context.Context, accountID, conversationID string, limit, offset int) ([]*Message, error) {
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found")
	}
	if conv.ParticipantA != accountID && conv.ParticipantB != accountID {
		return nil, fmt.Errorf("access denied")
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListMessages(ctx, conversationID, limit, offset)
}

func (s *Service) MarkAsRead(ctx context.Context, accountID, conversationID string) error {
	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("conversation not found")
	}
	if conv.ParticipantA != accountID && conv.ParticipantB != accountID {
		return fmt.Errorf("access denied")
	}
	if accountID == conv.ParticipantA {
		conv.UnreadCountA = 0
	} else {
		conv.UnreadCountB = 0
	}
	err = s.repo.UpdateConversation(ctx, conv)
	if err != nil {
		return fmt.Errorf("updating conversation: %w", err)
	}
	return s.repo.MarkMessagesAsRead(ctx, conversationID, accountID)
}
