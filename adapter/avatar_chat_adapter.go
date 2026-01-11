package adapter

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
)

type AvatarChatMessage struct {
	ID         string            `firestore:"id"`
	SenderType models.SenderType `firestore:"sender_type"`
	Message    string            `firestore:"message"`
	CreatedAt  time.Time         `firestore:"created_at"`
}

type AvatarChatAdapter interface {
	CreateAvatarChatMessage(ctx context.Context, userID string, avatarID string, message AvatarChatMessage) error
	GetAvatarChatMessages(ctx context.Context, userID string, avatarID string) ([]AvatarChatMessage, error)
}

type avatarChatAdapter struct {
	client *firestore.Client
}

func NewAvatarChatAdapter(client *firestore.Client) AvatarChatAdapter {
	return &avatarChatAdapter{
		client: client,
	}
}

func (a *avatarChatAdapter) CreateAvatarChatMessage(ctx context.Context, userID string, avatarID string, message AvatarChatMessage) error {
	col := a.client.Collection("avatar_chats").Doc(userID).Collection(avatarID)

	doc := col.Doc(message.ID)

	_, err := doc.Set(ctx, message)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}

func (a *avatarChatAdapter) GetAvatarChatMessages(ctx context.Context, userID string, avatarID string) ([]AvatarChatMessage, error) {
	query := a.client.Collection("avatar_chats").Doc(userID).Collection(avatarID).
		OrderBy("created_at", firestore.Asc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, utils.WrapError(err)
	}

	messages := make([]AvatarChatMessage, 0, len(docs))
	for _, doc := range docs {
		data := doc.Data()

		id, _ := data["id"].(string)
		senderType, _ := data["sender_type"].(string)
		message, _ := data["message"].(string)
		createdAt, _ := data["created_at"].(time.Time)

		messages = append(messages, AvatarChatMessage{
			ID:         id,
			SenderType: models.SenderType(senderType),
			Message:    message,
			CreatedAt:  createdAt,
		})
	}

	return messages, nil
}
