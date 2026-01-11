package response

// UnlockStatus アンロック状態情報
type UnlockStatus struct {
	ChatID          string `json:"chat_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	IsUnlocked      bool   `json:"is_unlocked" example:"false"`
	CurrentScore    int    `json:"current_score" example:"75"`
	UnlockThreshold int    `json:"unlock_threshold" example:"80"`
	RemainingScore  int    `json:"remaining_score" example:"5"` // アンロックまでに必要なスコア
	CanMatch        bool   `json:"can_match" example:"false"`   // マッチング可能かどうか
}

// GetUnlockStatusResponse アンロック状態取得レスポンス
type GetUnlockStatusResponse struct {
	Status UnlockStatus `json:"status"`
}

// Match マッチング情報
type Match struct {
	ID              string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FEV"`
	UserID          string `json:"user_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	PartnerID       string `json:"partner_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	ChatID          string `json:"chat_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	MatchedAt       string `json:"matched_at" example:"2024-01-01T12:00:00Z"`
	FinalScore      int    `json:"final_score" example:"85"`
	PartnerName     string `json:"partner_name" example:"佐藤花子"`
	PartnerImageURL string `json:"partner_image_url" example:"https://example.com/images/profile2.jpg"`
}

// MatchResponse マッチング成立レスポンス
type MatchResponse struct {
	Match   Match  `json:"match"`
	Message string `json:"message" example:"マッチングが成立しました！"`
}

// GetMatchesResponse マッチ一覧取得レスポンス
type GetMatchesResponse struct {
	Matches []Match `json:"matches"`
	Total   int     `json:"total" example:"5"`
}

// Suggestion 話題提案
type Suggestion struct {
	ID          string   `json:"id" example:"suggestion_123456"`
	Title       string   `json:"title" example:"共通の趣味について話す"`
	Description string   `json:"description" example:"お互いの趣味について詳しく聞いてみましょう"`
	Category    string   `json:"category" example:"interests"` // interests, hobbies, work, etc.
	Examples    []string `json:"examples,omitempty" example:"[\"どんな本が好きですか？\", \"最近読んだ本はありますか？\"]"`
}

// GetSuggestionsResponse 話題提案取得レスポンス
type GetSuggestionsResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// MatchMessage マッチ後のメッセージ情報
type MatchMessage struct {
	ID        string `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	MatchID   string `json:"match_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FEV"`
	SenderID  string `json:"sender_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	Content   string `json:"content" example:"こんにちは！マッチできて嬉しいです。"`
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

// GetMatchMessagesResponse マッチ後のメッセージ履歴取得レスポンス
type GetMatchMessagesResponse struct {
	Messages []MatchMessage `json:"messages"`
	Total    int            `json:"total" example:"20"`
	Limit    int            `json:"limit" example:"20"`
	Offset   int            `json:"offset" example:"0"`
}

// SendMatchMessageResponse マッチ後のメッセージ送信レスポンス
type SendMatchMessageResponse struct {
	Message MatchMessage `json:"message"`
}

// ReplyAssistResponse 返信アシストレスポンス
type ReplyAssistResponse struct {
	SuggestedReplies []string `json:"suggested_replies" example:"[\"ありがとうございます！こちらこそよろしくお願いします。\", \"こちらこそ！お話できるのを楽しみにしています。\"]"`
	Context          string   `json:"context,omitempty" example:"初めてのメッセージに対する返信"`
}
