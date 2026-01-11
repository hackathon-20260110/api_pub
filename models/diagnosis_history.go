package models

import "time"

type DiagnosisHistory struct {
	ID               string    `gorm:"primaryKey" json:"id"`                 // ULID
	UserID           string    `gorm:"not null" json:"user_id"`              // 診断実行者のID
	UserAvatarID     string    `gorm:"not null" json:"user_avatar_id"`       // 自分のAvatarのID
	TargetAvatarID   string    `gorm:"not null" json:"target_avatar_id"`     // 相手のAvatarのID
	ConversationData string    `gorm:"type:jsonb" json:"conversation_data"`  // Avatar同士の会話ログ（JSON形式）
	DiagnosisScore   int       `gorm:"not null" json:"diagnosis_score"`      // AI診断結果（1-5）
	AIAnalysisResult string    `gorm:"type:jsonb" json:"ai_analysis_result"` // AI分析の詳細結果（JSON形式）
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 外部キー関係
	User         User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserAvatar   Avatar `gorm:"foreignKey:UserAvatarID" json:"user_avatar,omitempty"`
	TargetAvatar Avatar `gorm:"foreignKey:TargetAvatarID" json:"target_avatar,omitempty"`
}
