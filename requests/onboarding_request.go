package requests

// SubmitBasicInfoRequest 基本情報送信リクエスト（オンボーディング開始前に必須）
type SubmitBasicInfoRequest struct {
	DisplayName        string `json:"display_name" example:"山田太郎" binding:"required"`
	Age                int    `json:"age" example:"25" binding:"required"`
	Gender             string `json:"gender" example:"male" binding:"required"` // male, female, other
	Bio                string `json:"bio,omitempty" example:"よろしくお願いします！"`
	ProfileImageBase64 string `json:"profile_image_base64,omitempty" example:"data:image/jpeg;base64,..."`
}

// StartOnboardingRequest オンボーディング開始リクエスト（チャット形式）
type StartOnboardingRequest struct {
	// 現在は空だが、将来的に初期設定などが必要になった場合に備える
	// 基本情報は既にSubmitBasicInfoで送信済み
}

// CompleteOnboardingRequest オンボーディング完了リクエスト
type CompleteOnboardingRequest struct {
	// 現在は空だが、将来的に最終確認などが必要になった場合に備える
}

// SendOnboardingMessageRequest オンボーディングチャットメッセージ送信リクエスト
type SendOnboardingMessageRequest struct {
	Content string `json:"content" example:"映画が好きです。" binding:"required"`
}
