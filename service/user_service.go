package service

import (
	"errors"
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type UserService struct {
	container *dig.Container
}

func NewUserService(container *dig.Container) *UserService {
	return &UserService{container: container}
}

func (s *UserService) GetUserByID(id string) (response.User, error) {
	var userAdapter adapter.UserAdapter
	if err := s.container.Invoke(func(adapter adapter.UserAdapter) error {
		userAdapter = adapter
		return nil
	}); err != nil {
		return response.User{}, utils.WrapError(err)
	}

	u, err := userAdapter.GetByID(id)
	if err != nil {
		return response.User{}, utils.WrapError(err)
	}

	r := response.NewUserResponse(u)

	return r, nil
}

// GetUserByIDOrNil ユーザーIDでユーザーを取得する。存在しない場合はnilを返す。
// 他のエラーが発生した場合はエラーを返す。
func (s *UserService) GetUserByIDOrNil(id string) (*response.User, error) {
	var userAdapter adapter.UserAdapter
	if err := s.container.Invoke(func(adapter adapter.UserAdapter) error {
		userAdapter = adapter
		return nil
	}); err != nil {
		return nil, err
	}

	u, err := userAdapter.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	r := response.NewUserResponse(u)

	return &r, nil
}

func (s *UserService) UpsertUser(userID string, args requests.CreateUserRequest) (response.User, error) {
	var userAdapter adapter.UserAdapter
	var r2Adapter adapter.R2Adapter
	if err := s.container.Invoke(func(ua adapter.UserAdapter, ra adapter.R2Adapter) error {
		userAdapter = ua
		r2Adapter = ra
		return nil
	}); err != nil {
		return response.User{}, err
	}

	imageData, err := utils.DecodeImageDataURI(args.ProfileImageBase64)
	if err != nil {
		return response.User{}, utils.WrapError(err)
	}

	objectKey := userID + imageData.Extension
	url, err := r2Adapter.UploadImage(imageData.Data, objectKey, imageData.ContentType)
	if err != nil {
		return response.User{}, utils.WrapError(err)
	}

	user := models.User{
		ID:                    userID,
		DisplayName:           args.DisplayName,
		Gender:                args.Gender,
		BirthDate:             args.BirthDate,
		Bio:                   args.Bio,
		ProfileImageURL:       url,
		IsOnboardingCompleted: false,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	var r response.User
	_, err = userAdapter.GetByID(userID)
	if err == utils.ErrorRecordNotFound {
		nu, err := userAdapter.Create(user)
		if err != nil {
			return response.User{}, utils.WrapError(err)
		}
		r = response.NewUserResponse(nu)
	} else if err != nil {
		return response.User{}, utils.WrapError(err)
	} else {
		nu, err := userAdapter.Update(user)
		if err != nil {
			return response.User{}, utils.WrapError(err)
		}
		r = response.NewUserResponse(nu)
	}

	return r, nil
}
