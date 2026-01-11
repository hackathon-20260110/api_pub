package adapter

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
)

type UserChatMessage struct {
	ID         string            `firestore:"id"`
	SenderID   string            `firestore:"sender_id"`
	SenderType models.SenderType `firestore:"sender_type"`
	Message    string            `firestore:"message"`
	CreatedAt  time.Time         `firestore:"created_at"`
}

type UserChatAdapter interface {
	CreateUserChatMessage(ctx context.Context, user1ID string, user2ID string, message UserChatMessage) error
	GetUserChatMessages(ctx context.Context, user1ID string, user2ID string) ([]UserChatMessage, error)
}

type userChatAdapter struct {
	client *firestore.Client
}

func NewUserChatAdapter(client *firestore.Client) UserChatAdapter {
	return &userChatAdapter{
		client: client,
	}
}

func normalizeUserIDs(userID1, userID2 string) (string, string) {
	if userID1 < userID2 {
		return userID1, userID2
	}
	return userID2, userID1
}

func (a *userChatAdapter) CreateUserChatMessage(ctx context.Context, user1ID string, user2ID string, message UserChatMessage) error {
	normalizedUser1, normalizedUser2 := normalizeUserIDs(user1ID, user2ID)

	col := a.client.Collection("user_chat").Doc(normalizedUser1).Collection(normalizedUser2)

	doc := col.Doc(message.ID)

	_, err := doc.Set(ctx, message)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}

func (a *userChatAdapter) GetUserChatMessages(ctx context.Context, user1ID string, user2ID string) ([]UserChatMessage, error) {
	normalizedUser1, normalizedUser2 := normalizeUserIDs(user1ID, user2ID)

	query := a.client.Collection("user_chat").Doc(normalizedUser1).Collection(normalizedUser2).
		OrderBy("created_at", firestore.Asc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, utils.WrapError(err)
	}

	messages := make([]UserChatMessage, 0, len(docs))
	for _, doc := range docs {
		data := doc.Data()

		id, _ := data["id"].(string)
		senderID, _ := data["sender_id"].(string)
		senderType, _ := data["sender_type"].(string)
		message, _ := data["message"].(string)
		createdAt, _ := data["created_at"].(time.Time)

		messages = append(messages, UserChatMessage{
			ID:         id,
			SenderID:   senderID,
			SenderType: models.SenderType(senderType),
			Message:    message,
			CreatedAt:  createdAt,
		})
	}

	return messages, nil
}
