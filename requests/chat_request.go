package requests

// CreateChatRequest 新しいチャット開始リクエスト
type CreateChatRequest struct {
	PartnerID string `json:"partner_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV" binding:"required"`
	// 最初のメッセージを同時に送信する場合
	InitialMessage string `json:"initial_message,omitempty" example:"こんにちは！"`
}

// SendMessageRequest メッセージ送信リクエスト
type SendMessageRequest struct {
	Content string `json:"content" example:"こんにちは！今日はいい天気ですね。" binding:"required"`
	// 将来的に画像やファイル添付などに対応する場合に拡張可能
}
