package requests

// UpdateUserRequest ユーザー更新リクエスト
type UpdateUserRequest struct {
	DisplayName     string `json:"display_name" example:"山田太郎"`
	Age             int    `json:"age" example:"25"`
	Gender          string `json:"gender" example:"male"`
	Bio             string `json:"bio" example:"よろしくお願いします！"`
	ProfileImageURL string `json:"profile_image_url" example:"https://example.com/images/profile.jpg"`
}
