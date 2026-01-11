package requests

type UpdateMatchingPointRequest struct {
	Points int `json:"points" validate:"required,min=1"`
}
