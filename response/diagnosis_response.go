package response

import "time"

// 診断実行結果のレスポンス
type DiagnosisResult struct {
	DiagnosisID      string                 `json:"diagnosis_id"`
	DiagnosisScore   int                    `json:"diagnosis_score"`   // 1-5
	PointsEarned     int                    `json:"points_earned"`     // 獲得ポイント
	AnalysisResult   map[string]interface{} `json:"analysis_result"`   // AI分析結果
	CanDirectChat    bool                   `json:"can_direct_chat"`   // 直接チャット可能フラグ
	CreatedAt        time.Time              `json:"created_at"`
}

// 診断履歴一覧のレスポンス
type DiagnosisHistory struct {
	ID               string    `json:"id"`
	UserAvatarName   string    `json:"user_avatar_name"`   // 自分のAvatar名
	TargetAvatarName string    `json:"target_avatar_name"` // 相手のAvatar名
	DiagnosisScore   int       `json:"diagnosis_score"`
	PointsEarned     int       `json:"points_earned"`
	CanDirectChat    bool      `json:"can_direct_chat"`
	CreatedAt        time.Time `json:"created_at"`
}

// 診断履歴一覧レスポンス
type DiagnosisHistoryResponse struct {
	Histories []DiagnosisHistory `json:"histories"`
}

// 診断詳細のレスポンス
type DiagnosisDetail struct {
	ID               string                 `json:"id"`
	UserAvatar       AvatarAI               `json:"user_avatar"`       // 自分のAvatar詳細
	TargetAvatar     AvatarAI               `json:"target_avatar"`     // 相手のAvatar詳細
	ConversationData map[string]interface{} `json:"conversation_data"` // 会話データ
	DiagnosisScore   int                    `json:"diagnosis_score"`
	PointsEarned     int                    `json:"points_earned"`
	AnalysisResult   map[string]interface{} `json:"analysis_result"` // AI分析詳細
	CanDirectChat    bool                   `json:"can_direct_chat"`
	CreatedAt        time.Time              `json:"created_at"`
}

// 診断詳細レスポンス
type DiagnosisDetailResponse struct {
	Diagnosis DiagnosisDetail `json:"diagnosis"`
}

// AvatarレスポンスのヘルパーFunction (既存のAvatarAI構造体を使用)
func NewAvatarResponse(avatar interface{}) AvatarAI {
	// modelsのAvatarからAvatarAIに変換
	// この実装は実際のAvatar modelsの構造に合わせて調整
	return AvatarAI{
		// 適切なマッピングロジックを後で実装
	}
}