package tests

import (
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/hackathon-20260110/api/tests/mock"
	"github.com/hackathon-20260110/api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
	"go.uber.org/mock/gomock"
)

// 1x1 PNG image bytes (smallest valid PNG)
var minimalPNGBytes = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk header
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1
	0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xDE, // bit depth, color type, etc
	0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, 0x54, // IDAT chunk header
	0x08, 0xD7, 0x63, 0xF8, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x05,
	0xFE, 0xD4, 0x7A, 0x02, // IDAT data + CRC
	0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, // IEND chunk
	0xAE, 0x42, 0x60, 0x82,
}

// 1x1 JPEG image bytes (smallest valid JPEG)
var minimalJPEGBytes = []byte{
	0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43,
	0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08, 0x07, 0x07, 0x07, 0x09,
	0x09, 0x08, 0x0A, 0x0C, 0x14, 0x0D, 0x0C, 0x0B, 0x0B, 0x0C, 0x19, 0x12,
	0x13, 0x0F, 0x14, 0x1D, 0x1A, 0x1F, 0x1E, 0x1D, 0x1A, 0x1C, 0x1C, 0x20,
	0x24, 0x2E, 0x27, 0x20, 0x22, 0x2C, 0x23, 0x1C, 0x1C, 0x28, 0x37, 0x29,
	0x2C, 0x30, 0x31, 0x34, 0x34, 0x34, 0x1F, 0x27, 0x39, 0x3D, 0x38, 0x32,
	0x3C, 0x2E, 0x33, 0x34, 0x32, 0xFF, 0xC0, 0x00, 0x0B, 0x08, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xFF, 0xC4, 0x00, 0x1F, 0x00, 0x00,
	0x01, 0x05, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
	0x09, 0x0A, 0x0B, 0xFF, 0xC4, 0x00, 0xB5, 0x10, 0x00, 0x02, 0x01, 0x03,
	0x03, 0x02, 0x04, 0x03, 0x05, 0x05, 0x04, 0x04, 0x00, 0x00, 0x01, 0x7D,
	0x01, 0x02, 0x03, 0x00, 0x04, 0x11, 0x05, 0x12, 0x21, 0x31, 0x41, 0x06,
	0x13, 0x51, 0x61, 0x07, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xA1, 0x08,
	0x23, 0x42, 0xB1, 0xC1, 0x15, 0x52, 0xD1, 0xF0, 0x24, 0x33, 0x62, 0x72,
	0x82, 0x09, 0x0A, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x25, 0x26, 0x27, 0x28,
	0x29, 0x2A, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x43, 0x44, 0x45,
	0x46, 0x47, 0x48, 0x49, 0x4A, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59,
	0x5A, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x73, 0x74, 0x75,
	0x76, 0x77, 0x78, 0x79, 0x7A, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89,
	0x8A, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9A, 0xA2, 0xA3,
	0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6,
	0xB7, 0xB8, 0xB9, 0xBA, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xC9,
	0xCA, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8, 0xD9, 0xDA, 0xE1, 0xE2,
	0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xF1, 0xF2, 0xF3, 0xF4,
	0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFF, 0xDA, 0x00, 0x08, 0x01, 0x01,
	0x00, 0x00, 0x3F, 0x00, 0xFB, 0xD5, 0xDB, 0x20, 0xA8, 0xF1, 0x7F, 0xFF,
	0xD9,
}

// Minimal HEIC ftyp box (not a full valid HEIC, but enough for detection)
var minimalHEICBytes = []byte{
	0x00, 0x00, 0x00, 0x18, // box size: 24 bytes
	0x66, 0x74, 0x79, 0x70, // box type: "ftyp"
	0x68, 0x65, 0x69, 0x63, // major brand: "heic"
	0x00, 0x00, 0x00, 0x00, // minor version
	0x68, 0x65, 0x69, 0x63, // compatible brand: "heic"
	0x6D, 0x69, 0x66, 0x31, // compatible brand: "mif1"
}

func buildDataURI(mimeType string, data []byte) string {
	return "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(data)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)

	now := time.Now()
	birthDate := time.Date(2000, 1, 15, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 2, 12, 0, 0, 0, time.UTC)

	expectedUser := models.User{
		ID:                    "user-123",
		DisplayName:           "テスト太郎",
		Gender:                "male",
		BirthDate:             birthDate,
		Bio:                   "よろしくお願いします",
		ProfileImageURL:       "https://example.com/image.jpg",
		IsOnboardingCompleted: true,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}

	mockUserAdapter.EXPECT().
		GetByID("user-123").
		Return(expectedUser, nil).
		Times(1)

	container := dig.New()
	err := container.Provide(func() adapter.UserAdapter {
		return mockUserAdapter
	})
	require.NoError(t, err)

	userService := service.NewUserService(container)
	result, err := userService.GetUserByID("user-123")

	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.DisplayName, result.DisplayName)
	assert.Equal(t, utils.CalculateAge(birthDate, now), result.Age)
	assert.Equal(t, expectedUser.ProfileImageURL, result.ProfileImageURL)
	assert.Equal(t, expectedUser.Bio, result.Bio)
	assert.Equal(t, expectedUser.IsOnboardingCompleted, result.OnboardingCompleted)
	assert.Equal(t, createdAt.Format(time.RFC3339), result.CreatedAt)
	assert.Equal(t, updatedAt.Format(time.RFC3339), result.UpdatedAt)
}

func TestUserService_GetUserByID_AdapterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)

	adapterErr := errors.New("user not found")
	mockUserAdapter.EXPECT().
		GetByID("nonexistent-user").
		Return(models.User{}, adapterErr).
		Times(1)

	container := dig.New()
	err := container.Provide(func() adapter.UserAdapter {
		return mockUserAdapter
	})
	require.NoError(t, err)

	userService := service.NewUserService(container)
	result, err := userService.GetUserByID("nonexistent-user")

	require.Error(t, err)
	assert.True(t, errors.Is(err, adapterErr))
	assert.Equal(t, response.User{}, result)
}

func TestUserService_GetUserByID_DIResolutionError(t *testing.T) {
	container := dig.New()

	userService := service.NewUserService(container)
	result, err := userService.GetUserByID("any-id")

	require.Error(t, err)
	assert.Equal(t, response.User{}, result)
}

func TestUserService_CreateUser_Success_PNG(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)
	mockR2Adapter := mock.NewMockR2Adapter(ctrl)

	now := time.Now()
	birthDate := time.Date(2000, 5, 20, 0, 0, 0, 0, time.UTC)
	userID := "firebase-uid-123"
	uploadedURL := "https://cdn.example.com/firebase-uid-123.png"

	pngDataURI := buildDataURI("image/png", minimalPNGBytes)

	request := requests.CreateUserRequest{
		DisplayName:        "新規ユーザー",
		Gender:             "female",
		BirthDate:          birthDate,
		Bio:                "初めまして！",
		ProfileImageBase64: pngDataURI,
	}

	mockR2Adapter.EXPECT().
		UploadImage(minimalPNGBytes, userID+".png", "image/png").
		Return(uploadedURL, nil).
		Times(1)

	mockUserAdapter.EXPECT().
		GetByID(userID).
		Return(models.User{}, utils.ErrorRecordNotFound).
		Times(1)

	mockUserAdapter.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(user models.User) (models.User, error) {
			assert.Equal(t, userID, user.ID)
			assert.Equal(t, request.DisplayName, user.DisplayName)
			assert.Equal(t, request.Gender, user.Gender)
			assert.Equal(t, request.BirthDate, user.BirthDate)
			assert.Equal(t, request.Bio, user.Bio)
			assert.Equal(t, uploadedURL, user.ProfileImageURL)
			assert.False(t, user.IsOnboardingCompleted)
			return user, nil
		}).
		Times(1)

	container := dig.New()
	require.NoError(t, container.Provide(func() adapter.UserAdapter { return mockUserAdapter }))
	require.NoError(t, container.Provide(func() adapter.R2Adapter { return mockR2Adapter }))

	userService := service.NewUserService(container)
	result, err := userService.UpsertUser(userID, request)

	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, request.DisplayName, result.DisplayName)
	assert.Equal(t, utils.CalculateAge(birthDate, now), result.Age)
	assert.Equal(t, uploadedURL, result.ProfileImageURL)
	assert.Equal(t, request.Bio, result.Bio)
	assert.False(t, result.OnboardingCompleted)
}

func TestUserService_CreateUser_Success_JPEG(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)
	mockR2Adapter := mock.NewMockR2Adapter(ctrl)

	birthDate := time.Date(2000, 5, 20, 0, 0, 0, 0, time.UTC)
	userID := "firebase-uid-456"
	uploadedURL := "https://cdn.example.com/firebase-uid-456.jpg"

	jpegDataURI := buildDataURI("image/jpeg", minimalJPEGBytes)

	request := requests.CreateUserRequest{
		DisplayName:        "JPEGユーザー",
		Gender:             "male",
		BirthDate:          birthDate,
		Bio:                "JPEG画像です",
		ProfileImageBase64: jpegDataURI,
	}

	mockR2Adapter.EXPECT().
		UploadImage(minimalJPEGBytes, userID+".jpg", "image/jpeg").
		Return(uploadedURL, nil).
		Times(1)

	mockUserAdapter.EXPECT().
		GetByID(userID).
		Return(models.User{}, utils.ErrorRecordNotFound).
		Times(1)

	mockUserAdapter.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(user models.User) (models.User, error) {
			return user, nil
		}).
		Times(1)

	container := dig.New()
	require.NoError(t, container.Provide(func() adapter.UserAdapter { return mockUserAdapter }))
	require.NoError(t, container.Provide(func() adapter.R2Adapter { return mockR2Adapter }))

	userService := service.NewUserService(container)
	result, err := userService.UpsertUser(userID, request)

	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, uploadedURL, result.ProfileImageURL)
}

func TestUserService_CreateUser_Success_HEIC(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)
	mockR2Adapter := mock.NewMockR2Adapter(ctrl)

	birthDate := time.Date(2000, 5, 20, 0, 0, 0, 0, time.UTC)
	userID := "firebase-uid-789"
	uploadedURL := "https://cdn.example.com/firebase-uid-789.heic"

	heicDataURI := buildDataURI("image/heic", minimalHEICBytes)

	request := requests.CreateUserRequest{
		DisplayName:        "HEICユーザー",
		Gender:             "female",
		BirthDate:          birthDate,
		Bio:                "HEIC画像です",
		ProfileImageBase64: heicDataURI,
	}

	mockR2Adapter.EXPECT().
		UploadImage(minimalHEICBytes, userID+".heic", "image/heic").
		Return(uploadedURL, nil).
		Times(1)

	mockUserAdapter.EXPECT().
		GetByID(userID).
		Return(models.User{}, utils.ErrorRecordNotFound).
		Times(1)

	mockUserAdapter.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(user models.User) (models.User, error) {
			return user, nil
		}).
		Times(1)

	container := dig.New()
	require.NoError(t, container.Provide(func() adapter.UserAdapter { return mockUserAdapter }))
	require.NoError(t, container.Provide(func() adapter.R2Adapter { return mockR2Adapter }))

	userService := service.NewUserService(container)
	result, err := userService.UpsertUser(userID, request)

	require.NoError(t, err)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, uploadedURL, result.ProfileImageURL)
}

func TestUserService_CreateUser_InvalidDataURI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)
	mockR2Adapter := mock.NewMockR2Adapter(ctrl)

	container := dig.New()
	require.NoError(t, container.Provide(func() adapter.UserAdapter { return mockUserAdapter }))
	require.NoError(t, container.Provide(func() adapter.R2Adapter { return mockR2Adapter }))

	userService := service.NewUserService(container)

	testCases := []struct {
		name               string
		profileImageBase64 string
		expectedErr        error
	}{
		{"raw base64 without data URI prefix", "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==", utils.ErrInvalidDataURIFormat},
		{"invalid data URI format", "invalid-data", utils.ErrInvalidDataURIFormat},
		{"unsupported MIME type", "data:image/webp;base64,UklGRlYAAABXRUJQVlA4IEoAAADQAQCdASoBAAEAAQAcJYgCdAEO/hOMAAD++O3u/v9W/xH/5seIjX5V/q/+LW8v/u21v///7N", utils.ErrUnsupportedMimeType},
		{"missing base64 marker", "data:image/png,iVBORw0KGgoAAAANSUhEUg==", utils.ErrInvalidDataURIFormat},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := requests.CreateUserRequest{
				DisplayName:        "テスト",
				ProfileImageBase64: tc.profileImageBase64,
			}

			result, err := userService.UpsertUser("user-id", request)

			require.Error(t, err)
			assert.True(t, errors.Is(err, tc.expectedErr))
			assert.Equal(t, response.User{}, result)
		})
	}
}

func TestUserService_CreateUser_R2UploadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)
	mockR2Adapter := mock.NewMockR2Adapter(ctrl)

	pngDataURI := buildDataURI("image/png", minimalPNGBytes)

	uploadErr := errors.New("failed to upload image")
	mockR2Adapter.EXPECT().
		UploadImage(minimalPNGBytes, "user-id.png", "image/png").
		Return("", uploadErr).
		Times(1)

	container := dig.New()
	require.NoError(t, container.Provide(func() adapter.UserAdapter { return mockUserAdapter }))
	require.NoError(t, container.Provide(func() adapter.R2Adapter { return mockR2Adapter }))

	userService := service.NewUserService(container)
	result, err := userService.UpsertUser("user-id", requests.CreateUserRequest{
		DisplayName:        "テスト",
		ProfileImageBase64: pngDataURI,
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, uploadErr))
	assert.Equal(t, response.User{}, result)
}

func TestUserService_CreateUser_UserAdapterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserAdapter := mock.NewMockUserAdapter(ctrl)
	mockR2Adapter := mock.NewMockR2Adapter(ctrl)

	pngDataURI := buildDataURI("image/png", minimalPNGBytes)

	mockR2Adapter.EXPECT().
		UploadImage(minimalPNGBytes, "user-id.png", "image/png").
		Return("https://cdn.example.com/image.png", nil).
		Times(1)

	mockUserAdapter.EXPECT().
		GetByID("user-id").
		Return(models.User{}, utils.ErrorRecordNotFound).
		Times(1)

	createErr := errors.New("database error")
	mockUserAdapter.EXPECT().
		Create(gomock.Any()).
		Return(models.User{}, createErr).
		Times(1)

	container := dig.New()
	require.NoError(t, container.Provide(func() adapter.UserAdapter { return mockUserAdapter }))
	require.NoError(t, container.Provide(func() adapter.R2Adapter { return mockR2Adapter }))

	userService := service.NewUserService(container)
	result, err := userService.UpsertUser("user-id", requests.CreateUserRequest{
		DisplayName:        "テスト",
		ProfileImageBase64: pngDataURI,
	})

	require.Error(t, err)
	assert.True(t, errors.Is(err, createErr))
	assert.Equal(t, response.User{}, result)
}

func TestUserService_CreateUser_DIResolutionError(t *testing.T) {
	container := dig.New()

	userService := service.NewUserService(container)
	result, err := userService.UpsertUser("user-id", requests.CreateUserRequest{})

	require.Error(t, err)
	assert.Equal(t, response.User{}, result)
}
