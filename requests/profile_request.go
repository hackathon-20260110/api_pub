package requests

type CreateUserInfoRequest struct {
	Key           string                `json:"key" binding:"required"`
	Value         string                `json:"value"`
	InfoType      string                `json:"info_type" binding:"required"`
	ImageBase64   string                `json:"image_base64,omitempty"`
	IsMission     bool                  `json:"is_mission"`
	MissionConfig *MissionConfigRequest `json:"mission_config,omitempty"`
}

type UpdateUserInfoRequest struct {
	Value         string                `json:"value"`
	ImageBase64   string                `json:"image_base64,omitempty"`
	IsMission     bool                  `json:"is_mission"`
	MissionConfig *MissionConfigRequest `json:"mission_config,omitempty"`
}

type MissionConfigRequest struct {
	ThresholdPoint  int    `json:"threshold_point" binding:"required,min=1"`
	UnlockCondition string `json:"unlock_condition,omitempty"`
}

