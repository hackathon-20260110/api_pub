package requests

// SendMatchMessageRequest マッチ後の本人とのメッセージ送信リクエスト
type SendMatchMessageRequest struct {
	Content string `json:"content" example:"こんにちは！マッチできて嬉しいです。" binding:"required"`
}

// ReplyAssistRequest 返信アシストリクエスト
type ReplyAssistRequest struct {
	PartnerMessage string `json:"partner_message" example:"こんにちは！こちらこそよろしくお願いします。" binding:"required"`
	// 会話の文脈を理解するための追加情報（必要に応じて）
	Context string `json:"context,omitempty" example:"初めてのメッセージ"`
}
