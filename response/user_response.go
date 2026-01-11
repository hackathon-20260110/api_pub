package response

// GetMeResponse 自分の情報取得レスポンス（マイページ用）
type GetMeResponse struct {
	IsRegistered bool  `json:"is_registered" example:"true"`
	User         *User `json:"user,omitempty"`
}

// UserDetail ユーザーの詳細情報（段階的公開）
type UserDetail struct {
	ID                 string   `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	DisplayName        string   `json:"display_name" example:"佐藤花子"`
	Age                int      `json:"age" example:"24"`
	Gender             string   `json:"gender" example:"female"`
	ProfileImageURL    string   `json:"profile_image_url" example:"https://example.com/images/profile2.jpg"`
	Bio                string   `json:"bio" example:"よろしくお願いします！"`
	Interests          []string `json:"interests,omitempty" example:"[\"読書\", \"映画\"]"`
	Location           string   `json:"location,omitempty" example:"東京都"`
	Occupation         string   `json:"occupation,omitempty" example:"エンジニア"`
	HasAvatarAI        bool     `json:"has_avatar_ai" example:"true"`
	AvatarAIAccessible bool     `json:"avatar_ai_accessible" example:"false"` // チャット開始済みかどうか
	CreatedAt          string   `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt          string   `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// GetUserDetailResponse 特定ユーザーの公開情報取得レスポンス
type GetUserDetailResponse struct {
	User UserDetail `json:"user"`
}

// GetUserAvatarAIResponse 特定ユーザーの分身AI取得レスポンス
type GetUserAvatarAIResponse struct {
	AvatarAI *AvatarAI `json:"avatar_ai"`
}
