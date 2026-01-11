package response

// Partner 相手候補の基本情報
type Partner struct {
	ID              string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	DisplayName     string `json:"display_name" example:"佐藤花子"`
	Age             int    `json:"age" example:"24"`
	Gender          string `json:"gender" example:"female"`
	ProfileImageURL string `json:"profile_image_url" example:"https://example.com/images/profile2.jpg"`
	Bio             string `json:"bio" example:"よろしくお願いします！"`
	// 基本情報のみを返す（詳細情報は段階的に公開）
}

// GetPartnersResponse 相手候補一覧取得レスポンス
type GetPartnersResponse struct {
	Partners []Partner `json:"partners"`
	Total    int       `json:"total" example:"50"`
	Limit    int       `json:"limit" example:"20"`
	Offset   int       `json:"offset" example:"0"`
}
