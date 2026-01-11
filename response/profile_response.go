package response

// UserProfileResponse ユーザープロフィール情報レスポンス
type UserProfileResponse struct {
	BasicInfo    BasicProfileInfo   `json:"basic_info"`
	UserInfoList []UserInfoResponse `json:"user_info_list"`
	Missions     []MissionResponse  `json:"missions"`
}

// BasicProfileInfo 基本プロフィール情報
type BasicProfileInfo struct {
	ID              string `json:"id"`
	DisplayName     string `json:"display_name"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_image_url"`
}

// UserInfoResponse プロフィール項目情報レスポンス
type UserInfoResponse struct {
	ID             string `json:"id"`
	Key            string `json:"key"`
	KeyDisplayName string `json:"key_display_name"`
	Value          string `json:"value"`
	InfoType       string `json:"info_type"`
	IsMission      bool   `json:"is_mission"`
	MissionID      string `json:"mission_id,omitempty"`
	IsUnlocked     bool   `json:"is_unlocked,omitempty"`
}

// MissionResponse ミッション情報レスポンス
type MissionResponse struct {
	ID                      string `json:"id"`
	UserInfoID              string `json:"user_info_id"`
	ThresholdPointCondition int    `json:"threshold_point_condition"`
	UnlockCondition         string `json:"unlock_condition,omitempty"`
	IsUnlocked              bool   `json:"is_unlocked"`
	UnlockedAt              string `json:"unlocked_at,omitempty"`
}

// PredefinedKeysResponse 事前定義項目一覧レスポンス
type PredefinedKeysResponse struct {
	Keys []PredefinedKeyInfo `json:"keys"`
}

// PredefinedKeyInfo 事前定義項目情報
type PredefinedKeyInfo struct {
	Key          string `json:"key"`
	DisplayName  string `json:"display_name"`
	InfoType     string `json:"info_type"`
	Placeholder  string `json:"placeholder"`
	CanBeMission bool   `json:"can_be_mission"`
}

// CreateUserInfoResponse プロフィール項目作成レスポンス
type CreateUserInfoResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
}

// UpdateUserInfoResponse プロフィール項目更新レスポンス
type UpdateUserInfoResponse struct {
	UserInfo UserInfoResponse `json:"user_info"`
}
