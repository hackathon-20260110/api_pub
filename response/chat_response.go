package response

// Chat チャット情報
type Chat struct {
	ID              string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	UserID          string `json:"user_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PartnerID       string `json:"partner_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	PartnerName     string `json:"partner_name" example:"佐藤花子"`
	PartnerImageURL string `json:"partner_image_url" example:"https://example.com/images/profile2.jpg"`
	LastMessage     string `json:"last_message,omitempty" example:"こんにちは！"`
	LastMessageAt   string `json:"last_message_at,omitempty" example:"2024-01-01T12:00:00Z"`
	UnreadCount     int    `json:"unread_count" example:"0"`
	MatchingScore   int    `json:"matching_score" example:"75"`
	IsMatched       bool   `json:"is_matched" example:"false"`
	CreatedAt       string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// CreateChatResponse 新しいチャット開始レスポンス
type CreateChatResponse struct {
	Chat Chat `json:"chat"`
}

// GetChatsResponse チャット一覧取得レスポンス
type GetChatsResponse struct {
	Chats []Chat `json:"chats"`
	Total int    `json:"total" example:"10"`
}

// GetChatDetailResponse 特定チャットの詳細取得レスポンス
type GetChatDetailResponse struct {
	Chat Chat `json:"chat"`
}

// Message メッセージ情報
type Message struct {
	ID         string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	ChatID     string `json:"chat_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	SenderID   string `json:"sender_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	SenderType string `json:"sender_type" example:"user"` // "user" or "avatar_ai"
	Content    string `json:"content" example:"こんにちは！今日はいい天気ですね。"`
	CreatedAt  string `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

// GetMessagesResponse メッセージ履歴取得レスポンス
type GetMessagesResponse struct {
	Messages []Message `json:"messages"`
	Total    int       `json:"total" example:"50"`
	Limit    int       `json:"limit" example:"20"`
	Offset   int       `json:"offset" example:"0"`
}

// SendMessageResponse メッセージ送信レスポンス
type SendMessageResponse struct {
	Message       Message  `json:"message"`
	AIResponse    *Message `json:"ai_response,omitempty"` // AIからの自動返信がある場合
	MatchingScore int      `json:"matching_score" example:"75"`
}

// ChatScore マッチングポイント情報
type ChatScore struct {
	ChatID          string         `json:"chat_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	CurrentScore    int            `json:"current_score" example:"75"`
	MaxScore        int            `json:"max_score" example:"100"`
	ScoreBreakdown  map[string]int `json:"score_breakdown,omitempty"`
	UnlockThreshold int            `json:"unlock_threshold" example:"80"` // マッチング成立に必要なスコア
	UpdatedAt       string         `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

// GetChatScoreResponse マッチングポイント取得レスポンス
type GetChatScoreResponse struct {
	Score ChatScore `json:"score"`
}
