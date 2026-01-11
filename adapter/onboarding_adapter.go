package adapter

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
	"gorm.io/gorm"
)

type OnboardingAdapter interface {
	CreateOnboardingChat(ctx context.Context, userID string, chat models.OnboardingChat) error
	GetOnboardingChats(ctx context.Context, userID string) ([]models.OnboardingChat, error)
}

type onboardingAdapter struct {
	client *firestore.Client
	db     *gorm.DB
}

func NewOnboardingAdapter(client *firestore.Client, db *gorm.DB) OnboardingAdapter {
	return &onboardingAdapter{
		client: client,
		db:     db,
	}
}

func (a *onboardingAdapter) CreateOnboardingChat(ctx context.Context, userID string, chat models.OnboardingChat) error {
	col := a.client.Collection("onboarding_chats").Doc(userID).Collection("chats")

	doc := col.Doc(chat.ID)

	_, err := doc.Set(ctx, chat)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}

func (a *onboardingAdapter) GetOnboardingChats(ctx context.Context, userID string) ([]models.OnboardingChat, error) {
	query := a.client.Collection("onboarding_chats").Doc(userID).Collection("chats").
		OrderBy("created_at", firestore.Asc)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, utils.WrapError(err)
	}

	chats := make([]models.OnboardingChat, 0, len(docs))
	for _, doc := range docs {
		data := doc.Data()

		id, _ := data["id"].(string)
		senderType, _ := data["sender_type"].(string)
		message, _ := data["message"].(string)
		isOnboardingCompleted, _ := data["is_onboarding_completed"].(bool)
		createdAt, _ := data["created_at"].(time.Time)

		chats = append(chats, models.OnboardingChat{
			ID:                    id,
			SenderType:            models.SenderType(senderType),
			Message:               message,
			IsOnboardingCompleted: isOnboardingCompleted,
			CreatedAt:             createdAt,
		})
	}

	return chats, nil
}
