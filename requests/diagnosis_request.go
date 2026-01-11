package requests

// 診断実行リクエスト
type ExecuteDiagnosisRequest struct {
	TargetAvatarID   string                 `json:"target_avatar_id" validate:"required" example:"01ARZ3NDEKTSV4RRFFQ69G5FBV"`
	ConversationData map[string]interface{} `json:"conversation_data" validate:"required"`
}

// 診断履歴取得のクエリパラメータ
type GetDiagnosisHistoryQuery struct {
	Limit  int `query:"limit" example:"20"`
	Offset int `query:"offset" example:"0"`
}