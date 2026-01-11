package requests

import "time"

type CreateUserRequest struct {
	DisplayName        string    `json:"display_name" example:"山田太郎"`
	Gender             string    `json:"gender" example:"male"`
	BirthDate          time.Time `json:"birth_date" example:"2000-01-01"`
	Bio                string    `json:"bio" example:"よろしくお願いします！"`
	ProfileImageBase64 string    `json:"profile_image_base64" example:"data:image/jpeg;base64,..."`
}
