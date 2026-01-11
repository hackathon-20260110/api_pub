package service

import (
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type AvatarService struct {
	container *dig.Container
}

func NewAvatarService(container *dig.Container) *AvatarService {
	return &AvatarService{container: container}
}

func (s *AvatarService) GetAvatarList(userID string) ([]response.AvatarWithRelation, error) {
	var avatarAdapter adapter.AvatarAdapter
	var userAdapter adapter.UserAdapter
	if err := s.container.Invoke(func(aa adapter.AvatarAdapter, ua adapter.UserAdapter) error {
		avatarAdapter = aa
		userAdapter = ua
		return nil
	}); err != nil {
		return nil, err
	}

	currentUser, err := userAdapter.GetByID(userID)
	if err != nil {
		return nil, err
	}

	avatars, err := avatarAdapter.GetOppositeGenderAvatars(currentUser.Gender)
	if err != nil {
		return nil, err
	}

	var result []response.AvatarWithRelation
	for _, avatar := range avatars {
		owner, err := userAdapter.GetByID(avatar.UserID)
		if err != nil {
			continue
		}

		avatarResponse := response.AvatarWithRelation{
			Avatar: response.Avatar{
				ID:                avatar.ID,
				UserID:            avatar.UserID,
				AvatarIconURL:     avatar.AvatarIconURL,
				Prompt:            avatar.Prompt,
				PersonalityTraits: avatar.PersonalityTraits,
				CreatedAt:         avatar.CreatedAt,
				UpdatedAt:         avatar.UpdatedAt,
			},
			UserDisplayName: owner.DisplayName,
			UserAge:         utils.CalculateAge(owner.BirthDate, time.Now()),
			UserBio:         owner.Bio,
		}

		relation, err := avatarAdapter.GetUserAvatarRelation(userID, avatar.ID)
		if err == nil {
			avatarResponse.Relation = &response.UserAvatarRelation{
				ID:            relation.ID,
				UserID:        relation.UserID,
				AvatarID:      relation.AvatarID,
				MatchingPoint: relation.MatchingPoint,
				CreatedAt:     relation.CreatedAt,
				UpdatedAt:     relation.UpdatedAt,
			}
		}

		result = append(result, avatarResponse)
	}

	return result, nil
}

func (s *AvatarService) UpdateMatchingPoint(userID, avatarID string, points int) (*response.UpdateMatchingPointResponse, error) {
	var avatarAdapter adapter.AvatarAdapter
	if err := s.container.Invoke(func(aa adapter.AvatarAdapter) error {
		avatarAdapter = aa
		return nil
	}); err != nil {
		return nil, err
	}

	relation, err := avatarAdapter.GetUserAvatarRelation(userID, avatarID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			relation = models.UserAvatarRelation{
				ID:            utils.GenerateULID(),
				UserID:        userID,
				AvatarID:      avatarID,
				MatchingPoint: points,
			}
			if err := avatarAdapter.CreateUserAvatarRelation(relation); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		relation.MatchingPoint += points
		if err := avatarAdapter.UpdateUserAvatarRelation(relation); err != nil {
			return nil, err
		}
	}

	return &response.UpdateMatchingPointResponse{
		ID:            relation.ID,
		UserID:        relation.UserID,
		AvatarID:      relation.AvatarID,
		MatchingPoint: relation.MatchingPoint,
		CreatedAt:     relation.CreatedAt,
		UpdatedAt:     relation.UpdatedAt,
	}, nil
}
