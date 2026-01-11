package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type ProfileService struct {
	container *dig.Container
}

func NewProfileService(container *dig.Container) *ProfileService {
	return &ProfileService{container: container}
}

// GetPredefinedInfoKeys 事前定義項目一覧を取得
func (s *ProfileService) GetPredefinedInfoKeys() []response.PredefinedKeyInfo {
	keys := make([]response.PredefinedKeyInfo, len(models.PredefinedInfoList))
	for i, info := range models.PredefinedInfoList {
		keys[i] = response.PredefinedKeyInfo{
			Key:          string(info.Key),
			DisplayName:  info.DisplayName,
			InfoType:     string(info.InfoType),
			Placeholder:  info.Placeholder,
			CanBeMission: info.CanBeMission,
		}
	}
	return keys
}

// CreateUserInfo プロフィール項目を作成
func (s *ProfileService) CreateUserInfo(ctx context.Context, userID string, req requests.CreateUserInfoRequest) (*response.UserInfoResponse, error) {
	var r2Adapter adapter.R2Adapter
	var db *gorm.DB

	if err := s.container.Invoke(func(ra adapter.R2Adapter, database *gorm.DB) error {
		r2Adapter = ra
		db = database
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	// バリデーション: 事前定義項目か確認
	predefinedKey := s.findPredefinedKey(req.Key)
	if predefinedKey == nil {
		return nil, fmt.Errorf("invalid key: %s is not a predefined key", req.Key)
	}

	// バリデーション: ミッション設定可能か確認
	if req.IsMission && !predefinedKey.CanBeMission {
		return nil, fmt.Errorf("key %s cannot be set as mission", req.Key)
	}

	// バリデーション: ミッション設定時はMissionConfigが必要
	if req.IsMission && req.MissionConfig == nil {
		return nil, fmt.Errorf("mission_config is required when is_mission is true")
	}

	// バリデーション: InfoTypeが一致しているか確認
	if string(predefinedKey.InfoType) != req.InfoType {
		return nil, fmt.Errorf("info_type mismatch: expected %s, got %s", predefinedKey.InfoType, req.InfoType)
	}

	var createdInfo models.UserInfo
	var createdMission *models.Mission

	// 1. 画像の場合はR2にアップロード（トランザクション外で実行）
	value := req.Value
	if req.InfoType == string(models.UserInfoTypeImage) && req.ImageBase64 != "" {
		imageData, err := utils.DecodeImageDataURI(req.ImageBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}

		objectKey := fmt.Sprintf("users/%s/info/%s%s", userID, utils.GenerateULID(), imageData.Extension)
		url, err := r2Adapter.UploadImage(imageData.Data, objectKey, imageData.ContentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
		value = url
	}

	// トランザクション開始
	err := db.Transaction(func(tx *gorm.DB) error {
		// 2. UserInfo作成
		userInfo := models.UserInfo{
			ID:              utils.GenerateULID(),
			UserID:          userID,
			InfoType:        models.UserInfoType(req.InfoType),
			Key:             req.Key,
			Value:           value,
			IsMissionReward: req.IsMission,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := tx.Create(&userInfo).Error; err != nil {
			return fmt.Errorf("failed to create user info: %w", err)
		}
		createdInfo = userInfo

		// 3. ミッション作成（必要な場合）
		if req.IsMission && req.MissionConfig != nil {
			mission := models.Mission{
				ID:                      utils.GenerateULID(),
				MissionOwnerUserID:      userID,
				UserInfoID:              createdInfo.ID,
				ThresholdPointCondition: &req.MissionConfig.ThresholdPoint,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			}

			if req.MissionConfig.UnlockCondition != "" {
				mission.UnlockCondition = &req.MissionConfig.UnlockCondition
			}

			if err := tx.Create(&mission).Error; err != nil {
				return fmt.Errorf("failed to create mission: %w", err)
			}
			createdMission = &mission
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// レスポンス構築
	return s.buildUserInfoResponse(createdInfo, createdMission, false), nil
}

// UpdateUserInfo プロフィール項目を更新
func (s *ProfileService) UpdateUserInfo(ctx context.Context, userID string, infoID string, req requests.UpdateUserInfoRequest) (*response.UserInfoResponse, error) {
	var profileAdapter adapter.ProfileAdapter
	var r2Adapter adapter.R2Adapter

	if err := s.container.Invoke(func(pa adapter.ProfileAdapter, ra adapter.R2Adapter) error {
		profileAdapter = pa
		r2Adapter = ra
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	// 既存のUserInfoを取得
	existingInfo, err := profileAdapter.GetUserInfoByID(infoID)
	if err != nil {
		return nil, fmt.Errorf("user info not found: %w", err)
	}

	// 権限チェック
	if existingInfo.UserID != userID {
		return nil, fmt.Errorf("unauthorized: user does not own this user info")
	}

	// 画像の場合はR2にアップロード
	value := req.Value
	if existingInfo.InfoType == models.UserInfoTypeImage && req.ImageBase64 != "" {
		imageData, err := utils.DecodeImageDataURI(req.ImageBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}

		objectKey := fmt.Sprintf("users/%s/info/%s%s", userID, utils.GenerateULID(), imageData.Extension)
		url, err := r2Adapter.UploadImage(imageData.Data, objectKey, imageData.ContentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
		value = url
	}

	// UserInfo更新
	existingInfo.Value = value
	existingInfo.IsMissionReward = req.IsMission
	existingInfo.UpdatedAt = time.Now()

	updatedInfo, err := profileAdapter.UpdateUserInfo(infoID, existingInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to update user info: %w", err)
	}

	// ミッション更新
	var mission *models.Mission
	if req.IsMission && req.MissionConfig != nil {
		existingMission, err := profileAdapter.GetMissionByUserInfoID(infoID)
		if err == nil {
			// 既存のミッションを更新
			existingMission.ThresholdPointCondition = &req.MissionConfig.ThresholdPoint
			if req.MissionConfig.UnlockCondition != "" {
				existingMission.UnlockCondition = &req.MissionConfig.UnlockCondition
			}
			existingMission.UpdatedAt = time.Now()
			updated, err := profileAdapter.UpdateMission(existingMission.ID, existingMission)
			if err != nil {
				return nil, fmt.Errorf("failed to update mission: %w", err)
			}
			mission = &updated
		} else {
			// 新規ミッション作成
			newMission := models.Mission{
				ID:                      utils.GenerateULID(),
				MissionOwnerUserID:      userID,
				UserInfoID:              infoID,
				ThresholdPointCondition: &req.MissionConfig.ThresholdPoint,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			}
			if req.MissionConfig.UnlockCondition != "" {
				newMission.UnlockCondition = &req.MissionConfig.UnlockCondition
			}
			created, err := profileAdapter.CreateMission(newMission)
			if err != nil {
				return nil, fmt.Errorf("failed to create mission: %w", err)
			}
			mission = &created
		}
	} else if req.IsMission == false {
		// ミッションを無効化する場合は削除
		if err := profileAdapter.DeleteMissionByUserInfoID(infoID); err != nil {
			return nil, fmt.Errorf("failed to delete mission: %w", err)
		}
	}

	return s.buildUserInfoResponse(updatedInfo, mission, false), nil
}

// DeleteUserInfo プロフィール項目を削除
func (s *ProfileService) DeleteUserInfo(ctx context.Context, userID string, infoID string) error {
	var profileAdapter adapter.ProfileAdapter

	if err := s.container.Invoke(func(pa adapter.ProfileAdapter) error {
		profileAdapter = pa
		return nil
	}); err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}

	// 既存のUserInfoを取得
	existingInfo, err := profileAdapter.GetUserInfoByID(infoID)
	if err != nil {
		return fmt.Errorf("user info not found: %w", err)
	}

	// 権限チェック
	if existingInfo.UserID != userID {
		return fmt.Errorf("unauthorized: user does not own this user info")
	}

	// 削除（関連するMissionも削除される）
	if err := profileAdapter.DeleteUserInfo(infoID); err != nil {
		return fmt.Errorf("failed to delete user info: %w", err)
	}

	return nil
}

// GetUserProfile ユーザーのプロフィール情報を取得
func (s *ProfileService) GetUserProfile(ctx context.Context, targetUserID string, viewerUserID string) (*response.UserProfileResponse, error) {
	var profileAdapter adapter.ProfileAdapter
	var userAdapter adapter.UserAdapter

	if err := s.container.Invoke(func(pa adapter.ProfileAdapter, ua adapter.UserAdapter) error {
		profileAdapter = pa
		userAdapter = ua
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	// 1. 基本情報取得
	user, err := userAdapter.GetByID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 2. UserInfo一覧取得
	userInfoList, err := profileAdapter.GetUserInfoByUserID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info list: %w", err)
	}

	// 3. Mission一覧取得
	missions, err := profileAdapter.GetMissionsByOwnerID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get missions: %w", err)
	}

	// 4. 自分のプロフィールの場合は全情報を返す
	isOwnProfile := targetUserID == viewerUserID

	// 5. 他ユーザー視点の場合は解禁状況を確認
	var unlockedMissions map[string]bool
	if !isOwnProfile {
		missionIDs := make([]string, len(missions))
		for i, m := range missions {
			missionIDs[i] = m.ID
		}
		unlocks, err := profileAdapter.GetMissionUnlocksByUserID(viewerUserID, missionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get mission unlocks: %w", err)
		}

		unlockedMissions = make(map[string]bool)
		for _, unlock := range unlocks {
			unlockedMissions[unlock.MissionID] = true
		}

		// マッチングポイントを取得して解禁判定
		matchingScore, err := profileAdapter.GetMatchingScore(ctx, viewerUserID, targetUserID)
		if err != nil {
			// マッチングポイント取得に失敗した場合は0として扱う
			matchingScore = 0
		}

		// 解禁条件を満たしているがまだ解禁されていないMissionを解禁
		for _, mission := range missions {
			if mission.ThresholdPointCondition != nil && *mission.ThresholdPointCondition <= matchingScore {
				if !unlockedMissions[mission.ID] {
					// 解禁レコードを作成
					unlock := models.MissionUnlock{
						ID:             utils.GenerateULID(),
						MissionID:      mission.ID,
						UnlockedUserID: viewerUserID,
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					}
					if _, err := profileAdapter.CreateMissionUnlock(unlock); err == nil {
						unlockedMissions[mission.ID] = true
					}
				}
			}
		}
	}

	// レスポンス構築
	return s.buildUserProfileResponse(user, userInfoList, missions, unlockedMissions, isOwnProfile), nil
}

// buildUserInfoResponse UserInfoレスポンスを構築
func (s *ProfileService) buildUserInfoResponse(userInfo models.UserInfo, mission *models.Mission, isUnlocked bool) *response.UserInfoResponse {
	resp := &response.UserInfoResponse{
		ID:             userInfo.ID,
		Key:            userInfo.Key,
		KeyDisplayName: s.getDisplayName(userInfo.Key),
		Value:          userInfo.Value,
		InfoType:       string(userInfo.InfoType),
		IsMission:      userInfo.IsMissionReward,
		IsUnlocked:     isUnlocked,
	}

	if mission != nil {
		resp.MissionID = mission.ID
	}

	return resp
}

// buildUserProfileResponse ユーザープロフィールレスポンスを構築
func (s *ProfileService) buildUserProfileResponse(user models.User, userInfoList []models.UserInfo, missions []models.Mission, unlockedMissions map[string]bool, isOwnProfile bool) *response.UserProfileResponse {
	// MissionIDからUserInfoIDへのマッピングを作成
	missionMap := make(map[string]models.Mission)
	for _, mission := range missions {
		missionMap[mission.UserInfoID] = mission
	}

	// UserInfoレスポンスリストを構築
	userInfoResponses := make([]response.UserInfoResponse, 0, len(userInfoList))
	for _, userInfo := range userInfoList {
		mission, hasMission := missionMap[userInfo.ID]
		isUnlocked := false

		if hasMission {
			if isOwnProfile {
				isUnlocked = true
			} else {
				isUnlocked = unlockedMissions[mission.ID]
			}

			// 解禁されていない場合は「？？？」を表示
			if !isUnlocked {
				userInfo.Value = "？？？"
			}
		} else if isOwnProfile {
			isUnlocked = true
		}

		resp := s.buildUserInfoResponse(userInfo, &mission, isUnlocked)
		if !hasMission {
			resp.MissionID = ""
		}
		userInfoResponses = append(userInfoResponses, *resp)
	}

	// Missionレスポンスリストを構築
	missionResponses := make([]response.MissionResponse, 0, len(missions))
	for _, mission := range missions {
		isUnlocked := false
		if isOwnProfile {
			isUnlocked = true
		} else {
			isUnlocked = unlockedMissions[mission.ID]
		}

		threshold := 0
		if mission.ThresholdPointCondition != nil {
			threshold = *mission.ThresholdPointCondition
		}

		unlockCondition := ""
		if mission.UnlockCondition != nil {
			unlockCondition = *mission.UnlockCondition
		}

		missionResp := response.MissionResponse{
			ID:                      mission.ID,
			UserInfoID:              mission.UserInfoID,
			ThresholdPointCondition: threshold,
			UnlockCondition:         unlockCondition,
			IsUnlocked:              isUnlocked,
		}

		if isUnlocked {
			missionResp.UnlockedAt = time.Now().Format(time.RFC3339)
		}

		missionResponses = append(missionResponses, missionResp)
	}

	// 基本情報を構築
	age := utils.CalculateAge(user.BirthDate, time.Now())
	basicInfo := response.BasicProfileInfo{
		ID:              user.ID,
		DisplayName:     user.DisplayName,
		Age:             age,
		Gender:          user.Gender,
		Bio:             user.Bio,
		ProfileImageURL: user.ProfileImageURL,
	}

	return &response.UserProfileResponse{
		BasicInfo:    basicInfo,
		UserInfoList: userInfoResponses,
		Missions:     missionResponses,
	}
}

// getDisplayName キーから表示名を取得
func (s *ProfileService) getDisplayName(key string) string {
	for _, info := range models.PredefinedInfoList {
		if string(info.Key) == key {
			return info.DisplayName
		}
	}
	return key
}

// findPredefinedKey キーから事前定義項目を検索
func (s *ProfileService) findPredefinedKey(key string) *models.PredefinedInfoMetadata {
	for _, info := range models.PredefinedInfoList {
		if string(info.Key) == key {
			return &info
		}
	}
	return nil
}
