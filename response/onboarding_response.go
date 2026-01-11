package response

import "github.com/hackathon-20260110/api/models"

// StartOnboardingResponse オンボーディング開始レスポンス（チャット形式）
type StartOnboardingResponse struct {
	ChatRoomID string `json:"chat_room_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Message    string `json:"message" example:"オンボーディングチャットを開始しました。Firestoreからメッセージを取得してください。"`
}

// AvatarAI 分身AI情報
type AvatarAI struct {
	ID          string            `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FDV"`
	UserID      string            `json:"user_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Name        string            `json:"name" example:"私の分身AI"`
	Personality map[string]string `json:"personality" example:"{\"tone\": \"friendly\", \"style\": \"casual\"}"`
	Bio         string            `json:"bio" example:"あなたの性格や好みを反映した分身AIです"`
	CreatedAt   string            `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string            `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// GetAvatarAIResponse 分身AI情報取得レスポンス
type GetAvatarAIResponse struct {
	AvatarAI *AvatarAI `json:"avatar_ai"`
}

// OnboardingMessage オンボーディングチャットメッセージ情報
type OnboardingMessage struct {
	ID         string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	ChatRoomID string `json:"chat_room_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	SenderType string `json:"sender_type" example:"user"` // "user" or "avatar_ai"
	Content    string `json:"content" example:"こんにちは！今日はいい天気ですね。"`
	Category   string `json:"category,omitempty" example:"interests"` // 質問カテゴリ（オプション）
	CreatedAt  string `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

// SendOnboardingMessageResponse オンボーディングチャットメッセージ送信レスポンス
type SendOnboardingMessageResponse struct {
	Message string `json:"message" example:"メッセージを受け付けました"`
}

// FinishOnboardingResponse オンボーディング完了レスポンス
type FinishOnboardingResponse struct {
	Message string `json:"message" example:"オンボーディングチャット完了成功"`
}

type OnboardingCompleteResponse struct {
	UserInfos []models.UserInfo `json:"user_infos"`
	Avatar    models.Avatar     `json:"avatar"`
}
