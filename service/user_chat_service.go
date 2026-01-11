package service

import (
	"context"
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
)

type UserChatService struct {
	container *dig.Container
}

func NewUserChatService(container *dig.Container) *UserChatService {
	return &UserChatService{container: container}
}

type MatchedUserInfo struct {
	UserID          string
	DisplayName     string
	Gender          string
	Bio             string
	ProfileImageURL string
	MatchedAt       time.Time
}

func (s *UserChatService) GetMatchedUsers(ctx context.Context, userID string) ([]MatchedUserInfo, error) {
	var matchingAdapter adapter.MatchingAdapter
	var userAdapter adapter.UserAdapter

	if err := s.container.Invoke(func(
		ma adapter.MatchingAdapter,
		ua adapter.UserAdapter,
	) error {
		matchingAdapter = ma
		userAdapter = ua
		return nil
	}); err != nil {
		return nil, utils.WrapError(err)
	}

	matchings, err := matchingAdapter.GetMatchingsByUserID(userID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	matchedUsers := make([]MatchedUserInfo, 0, len(matchings))
	for _, matching := range matchings {
		var partnerID string
		if matching.User1ID == userID {
			partnerID = matching.User2ID
		} else {
			partnerID = matching.User1ID
		}

		partner, err := userAdapter.GetByID(partnerID)
		if err != nil {
			continue
		}

		matchedUsers = append(matchedUsers, MatchedUserInfo{
			UserID:          partner.ID,
			DisplayName:     partner.DisplayName,
			Gender:          partner.Gender,
			Bio:             partner.Bio,
			ProfileImageURL: partner.ProfileImageURL,
			MatchedAt:       matching.CreatedAt,
		})
	}

	return matchedUsers, nil
}

func (s *UserChatService) SendMessage(ctx context.Context, senderID string, partnerID string, content string) (*adapter.UserChatMessage, error) {
	var userChatAdapter adapter.UserChatAdapter
	var matchingAdapter adapter.MatchingAdapter

	if err := s.container.Invoke(func(
		uca adapter.UserChatAdapter,
		ma adapter.MatchingAdapter,
	) error {
		userChatAdapter = uca
		matchingAdapter = ma
		return nil
	}); err != nil {
		return nil, utils.WrapError(err)
	}

	_, err := matchingAdapter.GetMatchingByUsers(senderID, partnerID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	message := adapter.UserChatMessage{
		ID:         utils.GenerateULID(),
		SenderID:   senderID,
		SenderType: models.SenderTypeUser,
		Message:    content,
		CreatedAt:  time.Now(),
	}

	if err := userChatAdapter.CreateUserChatMessage(ctx, senderID, partnerID, message); err != nil {
		return nil, utils.WrapError(err)
	}

	return &message, nil
}

func (s *UserChatService) GetMessages(ctx context.Context, userID string, partnerID string) ([]adapter.UserChatMessage, error) {
	var userChatAdapter adapter.UserChatAdapter
	var matchingAdapter adapter.MatchingAdapter

	if err := s.container.Invoke(func(
		uca adapter.UserChatAdapter,
		ma adapter.MatchingAdapter,
	) error {
		userChatAdapter = uca
		matchingAdapter = ma
		return nil
	}); err != nil {
		return nil, utils.WrapError(err)
	}

	_, err := matchingAdapter.GetMatchingByUsers(userID, partnerID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	messages, err := userChatAdapter.GetUserChatMessages(ctx, userID, partnerID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	return messages, nil
}
