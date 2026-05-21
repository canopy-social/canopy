package dm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sumi-devs/canopy-social/canopy/internal/accounts"
	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
)

type mockAccountLoader struct {
	accounts map[string]*accounts.Account
}

func (m *mockAccountLoader) GetByID(ctx context.Context, id string) (*accounts.Account, error) {
	acc, ok := m.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}

type mockRepository struct {
	conversations map[string]*Conversation
	messages      map[string]*Message
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		conversations: make(map[string]*Conversation),
		messages:      make(map[string]*Message),
	}
}

func (m *mockRepository) CreateConversation(ctx context.Context, conv *Conversation) (*Conversation, error) {
	m.conversations[conv.ID] = conv
	return conv, nil
}

func (m *mockRepository) GetConversationByID(ctx context.Context, id string) (*Conversation, error) {
	conv, ok := m.conversations[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return conv, nil
}

func (m *mockRepository) GetConversationByParticipants(ctx context.Context, p1, p2 string) (*Conversation, error) {
	for _, c := range m.conversations {
		if c.ParticipantA == p1 && c.ParticipantB == p2 {
			return c, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepository) ListConversations(ctx context.Context, accountID string, limit, offset int) ([]*Conversation, error) {
	var res []*Conversation
	for _, c := range m.conversations {
		if c.ParticipantA == accountID || c.ParticipantB == accountID {
			res = append(res, c)
		}
	}
	return res, nil
}

func (m *mockRepository) UpdateConversation(ctx context.Context, conv *Conversation) error {
	m.conversations[conv.ID] = conv
	return nil
}

func (m *mockRepository) CreateMessage(ctx context.Context, msg *Message) (*Message, error) {
	m.messages[msg.ID] = msg
	return msg, nil
}

func (m *mockRepository) GetMessageByID(ctx context.Context, id string) (*Message, error) {
	msg, ok := m.messages[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return msg, nil
}

func (m *mockRepository) ListMessages(ctx context.Context, conversationID string, limit, offset int) ([]*Message, error) {
	var res []*Message
	for _, msg := range m.messages {
		if msg.ConversationID == conversationID {
			res = append(res, msg)
		}
	}
	return res, nil
}

func (m *mockRepository) MarkMessagesAsRead(ctx context.Context, conversationID, readerID string) error {
	now := time.Now()
	for _, msg := range m.messages {
		if msg.ConversationID == conversationID && msg.SenderID != readerID {
			msg.ReadAt = &now
		}
	}
	return nil
}

func TestDM_SendDM(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{
		accounts: map[string]*accounts.Account{
			"alice": {ID: "alice"},
			"bob":   {ID: "bob"},
		},
	}
	svc := NewService(repo, al)
	ctx := context.Background()
	msg, err := svc.SendDM(ctx, "alice", "bob", "hello bob")
	if err != nil {
		t.Fatalf("failed to send DM: %v", err)
	}
	if msg.Content != "hello bob" {
		t.Errorf("expected content 'hello bob', got '%s'", msg.Content)
	}
	if len(repo.conversations) != 1 {
		t.Fatalf("expected 1 conversation, got %d", len(repo.conversations))
	}
	var conv *Conversation
	for _, c := range repo.conversations {
		conv = c
	}
	if conv.ParticipantA != "alice" || conv.ParticipantB != "bob" {
		t.Errorf("expected participants alice and bob, got %s and %s", conv.ParticipantA, conv.ParticipantB)
	}
	if conv.UnreadCountB != 1 {
		t.Errorf("expected bob unread count to be 1, got %d", conv.UnreadCountB)
	}
	msg2, err := svc.SendDM(ctx, "bob", "alice", "hey alice")
	if err != nil {
		t.Fatalf("failed to reply: %v", err)
	}
	if len(repo.conversations) != 1 {
		t.Fatalf("expected conversation count to still be 1, got %d", len(repo.conversations))
	}
	if conv.UnreadCountA != 1 {
		t.Errorf("expected alice unread count to be 1, got %d", conv.UnreadCountA)
	}
	if msg2.ConversationID != msg.ConversationID {
		t.Errorf("expected same conversation ID, got %s vs %s", msg2.ConversationID, msg.ConversationID)
	}
}

func TestDM_MarkAsRead(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{
		accounts: map[string]*accounts.Account{
			"alice": {ID: "alice"},
			"bob":   {ID: "bob"},
		},
	}
	svc := NewService(repo, al)
	ctx := context.Background()
	_, err := svc.SendDM(ctx, "alice", "bob", "hello bob")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	var conv *Conversation
	for _, c := range repo.conversations {
		conv = c
	}
	if conv.UnreadCountB != 1 {
		t.Errorf("expected unread count 1, got %d", conv.UnreadCountB)
	}
	err = svc.MarkAsRead(ctx, "bob", conv.ID)
	if err != nil {
		t.Fatalf("failed mark as read: %v", err)
	}
	if conv.UnreadCountB != 0 {
		t.Errorf("expected reset unread count, got %d", conv.UnreadCountB)
	}
	var msg *Message
	for _, m := range repo.messages {
		msg = m
	}
	if msg.ReadAt == nil {
		t.Error("expected message to be marked read")
	}
}

func TestDM_Handler(t *testing.T) {
	repo := newMockRepository()
	al := &mockAccountLoader{
		accounts: map[string]*accounts.Account{
			"alice": {ID: "alice"},
			"bob":   {ID: "bob"},
		},
	}
	svc := NewService(repo, al)
	h := NewHandler(svc)
	jwtMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), auth.ContextKeyAccountID, "alice")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	r := chi.NewRouter()
	h.RegisterRoutes(r, jwtMiddleware)
	server := httptest.NewServer(r)
	defer server.Close()
	payload := `{"recipient_id":"bob","content":"hey bob!"}`
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/direct_messages", strings.NewReader(payload))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed post request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	var createdMsg Message
	json.NewDecoder(resp.Body).Decode(&createdMsg)
	if createdMsg.Content != "hey bob!" {
		t.Errorf("unexpected message response content: %s", createdMsg.Content)
	}
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/direct_messages/conversations", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed get conversations: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var list []*Conversation
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 1 {
		t.Errorf("expected 1 conversation, got %d", len(list))
	}
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/direct_messages/conversations/"+createdMsg.ConversationID+"/messages", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed get messages: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var msgs []*Message
	json.NewDecoder(resp.Body).Decode(&msgs)
	if len(msgs) != 1 {
		t.Errorf("expected 1 message, got %d", len(msgs))
	}
}
