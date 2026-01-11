package service

import (
	"context"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
)

type NotificationService struct {
	container *dig.Container
}

func NewNotificationService(container *dig.Container) *NotificationService {
	return &NotificationService{container: container}
}

func (s *NotificationService) MarkAsRead(ctx context.Context, userID string, notificationID string) error {
	var notificationAdapter adapter.NotificationAdapter

	if err := s.container.Invoke(func(na adapter.NotificationAdapter) error {
		notificationAdapter = na
		return nil
	}); err != nil {
		return utils.WrapError(err)
	}

	if err := notificationAdapter.MarkAsRead(ctx, userID, notificationID); err != nil {
		return utils.WrapError(err)
	}

	return nil
}
