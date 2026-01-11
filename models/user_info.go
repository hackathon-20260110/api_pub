package models

import "time"

type UserInfoType string

const (
	UserInfoTypeText  UserInfoType = "text"
	UserInfoTypeImage UserInfoType = "image"
)

// PredefinedInfoKey 事前定義されたプロフィール項目のキー
type PredefinedInfoKey string

const (
	// テキスト項目
	InfoKeyHobby        PredefinedInfoKey = "hobby"
	InfoKeyFavoriteFood PredefinedInfoKey = "favorite_food"
	InfoKeyBloodType    PredefinedInfoKey = "blood_type"
	InfoKeyWork         PredefinedInfoKey = "work"
	InfoKeyEducation    PredefinedInfoKey = "education"
	InfoKeyHometown     PredefinedInfoKey = "hometown"
	InfoKeyLanguages    PredefinedInfoKey = "languages"
	InfoKeyPet          PredefinedInfoKey = "pet"
	InfoKeyDrinking     PredefinedInfoKey = "drinking"
	InfoKeySmoking      PredefinedInfoKey = "smoking"

	// 画像項目（2枚目以降）
	InfoKeySubImage1 PredefinedInfoKey = "sub_image_1"
	InfoKeySubImage2 PredefinedInfoKey = "sub_image_2"
	InfoKeySubImage3 PredefinedInfoKey = "sub_image_3"
)

// PredefinedInfoMetadata 事前定義項目のメタデータ
type PredefinedInfoMetadata struct {
	Key          PredefinedInfoKey
	DisplayName  string
	InfoType     UserInfoType
	Placeholder  string
	CanBeMission bool
}

// PredefinedInfoList 事前定義項目のリスト
var PredefinedInfoList = []PredefinedInfoMetadata{
	{
		Key:          InfoKeyHobby,
		DisplayName:  "趣味",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 読書、映画鑑賞、ランニング",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyFavoriteFood,
		DisplayName:  "好きな食べ物",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: イタリアン、和食",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyBloodType,
		DisplayName:  "血液型",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: A型、B型、O型、AB型",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyWork,
		DisplayName:  "仕事",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: エンジニア、デザイナー",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyEducation,
		DisplayName:  "学歴",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 大学卒業、大学院卒業",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyHometown,
		DisplayName:  "出身地",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 東京都、大阪府",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyLanguages,
		DisplayName:  "言語",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 日本語、英語、中国語",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyPet,
		DisplayName:  "ペット",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 犬、猫、なし",
		CanBeMission: true,
	},
	{
		Key:          InfoKeyDrinking,
		DisplayName:  "飲酒",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 飲む、飲まない、たまに",
		CanBeMission: true,
	},
	{
		Key:          InfoKeySmoking,
		DisplayName:  "喫煙",
		InfoType:     UserInfoTypeText,
		Placeholder:  "例: 吸う、吸わない",
		CanBeMission: true,
	},
	{
		Key:          InfoKeySubImage1,
		DisplayName:  "サブ画像1",
		InfoType:     UserInfoTypeImage,
		Placeholder:  "",
		CanBeMission: true,
	},
	{
		Key:          InfoKeySubImage2,
		DisplayName:  "サブ画像2",
		InfoType:     UserInfoTypeImage,
		Placeholder:  "",
		CanBeMission: true,
	},
	{
		Key:          InfoKeySubImage3,
		DisplayName:  "サブ画像3",
		InfoType:     UserInfoTypeImage,
		Placeholder:  "",
		CanBeMission: true,
	},
}

type UserInfo struct {
	ID       string       `gorm:"primaryKey" json:"id"`
	UserID   string       `json:"user_id" gorm:"not null"`
	InfoType UserInfoType `json:"type" gorm:"not null"`
	// テキストの場合 Key: 項目名、Value: ユーザーが入力した値
	// 画像の場合 Key: 画像タイトル、Value: 画像のURL
	Key             string    `json:"key" gorm:"not null"`
	Value           string    `json:"value" gorm:"not null"`
	IsMissionReward bool      `json:"is_mission_reward" gorm:"default:false"` // ミッション報酬になっているかどうか
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
