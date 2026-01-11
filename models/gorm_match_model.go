package models

import "time"

type Match struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	ChatID     string    `json:"chat_id"`
	UserAID    string    `json:"user_a_id"`
	UserBID    string    `json:"user_b_id"`
	FinalScore int       `json:"final_score"`
	MatchedAt  time.Time `gorm:"autoCreateTime" json:"matched_at"`
}

func (Match) TableName() string {
	return "matches"
}
