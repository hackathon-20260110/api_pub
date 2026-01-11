package adapter

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
)

type FirestoreNotification struct {
	ID        string    `firestore:"id"`
	UserID    string    `firestore:"user_id"`
	Title     string    `firestore:"title"`
	Message   string    `firestore:"message"`
	HasRead   bool      `firestore:"has_read"`
	CreatedAt time.Time `firestore:"created_at"`
}

type NotificationAdapter interface {
	CreateNotification(ctx context.Context, userID string, notification models.Notification) error
	MarkAsRead(ctx context.Context, userID string, notificationID string) error
}

type notificationAdapter struct {
	client *firestore.Client
}

func NewNotificationAdapter(client *firestore.Client) NotificationAdapter {
	return &notificationAdapter{
		client: client,
	}
}

func (a *notificationAdapter) CreateNotification(ctx context.Context, userID string, notification models.Notification) error {
	col := a.client.Collection("notifications").Doc(userID).Collection("notification")
	doc := col.Doc(notification.ID)

	firestoreNotification := FirestoreNotification{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		HasRead:   false,
		CreatedAt: notification.CreatedAt,
	}

	_, err := doc.Set(ctx, firestoreNotification)
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}

func (a *notificationAdapter) MarkAsRead(ctx context.Context, userID string, notificationID string) error {
	doc := a.client.Collection("notifications").Doc(userID).Collection("notification").Doc(notificationID)

	_, err := doc.Update(ctx, []firestore.Update{
		{Path: "has_read", Value: true},
	})
	if err != nil {
		return utils.WrapError(err)
	}

	return nil
}
