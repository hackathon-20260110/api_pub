package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/utils"
	"go.uber.org/dig"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

type AvatarChatService struct {
	container *dig.Container
}

func NewAvatarChatService(container *dig.Container) *AvatarChatService {
	return &AvatarChatService{container: container}
}

type LLMChatResponse struct {
	Message     string `json:"message"`
	PointChange int    `json:"point_change"`
	Reason      string `json:"reason"`
}

type SendMessageResult struct {
	AvatarResponse   adapter.AvatarChatMessage
	MatchingPoint    int
	PointChange      int
	IsMatched        bool
	UnlockedMissions []UnlockedMissionInfo
}

type UnlockedMissionInfo struct {
	MissionID  string
	UserInfoID string
	Key        string
	Value      string
}

func (s *AvatarChatService) SendMessage(ctx context.Context, userID string, avatarID string, content string) (*SendMessageResult, error) {
	var avatarChatAdapter adapter.AvatarChatAdapter
	var avatarAdapter adapter.AvatarAdapter
	var userAdapter adapter.UserAdapter
	var userInfoAdapter adapter.UserInfoAdapter
	var missionAdapter adapter.MissionAdapter
	var matchingAdapter adapter.MatchingAdapter
	var llmAdapter adapter.LLMAdapter
	var notificationAdapter adapter.NotificationAdapter

	if err := s.container.Invoke(func(
		aca adapter.AvatarChatAdapter,
		aa adapter.AvatarAdapter,
		ua adapter.UserAdapter,
		uia adapter.UserInfoAdapter,
		ma adapter.MissionAdapter,
		mta adapter.MatchingAdapter,
		la adapter.LLMAdapter,
		na adapter.NotificationAdapter,
	) error {
		avatarChatAdapter = aca
		avatarAdapter = aa
		userAdapter = ua
		userInfoAdapter = uia
		missionAdapter = ma
		matchingAdapter = mta
		llmAdapter = la
		notificationAdapter = na
		return nil
	}); err != nil {
		return nil, utils.WrapError(err)
	}

	avatar, err := avatarAdapter.GetByID(avatarID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	avatarOwnerUser, err := userAdapter.GetByID(avatar.UserID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	avatarOwnerUserInfos, err := userInfoAdapter.GetByUserID(avatar.UserID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	relation, err := avatarAdapter.GetUserAvatarRelation(userID, avatarID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			relation = models.UserAvatarRelation{
				ID:            utils.GenerateULID(),
				UserID:        userID,
				AvatarID:      avatarID,
				MatchingPoint: 0,
			}
			if err := avatarAdapter.CreateUserAvatarRelation(relation); err != nil {
				return nil, utils.WrapError(err)
			}
		} else {
			return nil, utils.WrapError(err)
		}
	}

	existingMatching, _ := matchingAdapter.GetMatchingByUsers(userID, avatar.UserID)
	isAlreadyMatched := existingMatching != nil

	userMessage := adapter.AvatarChatMessage{
		ID:         utils.GenerateULID(),
		SenderType: models.SenderTypeUser,
		Message:    content,
		CreatedAt:  time.Now(),
	}

	if err := avatarChatAdapter.CreateAvatarChatMessage(ctx, userID, avatarID, userMessage); err != nil {
		return nil, utils.WrapError(err)
	}

	chatHistory, err := avatarChatAdapter.GetAvatarChatMessages(ctx, userID, avatarID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	llmResponse, err := s.generateAvatarResponse(avatarOwnerUser, avatarOwnerUserInfos, chatHistory, llmAdapter)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	avatarResponse := adapter.AvatarChatMessage{
		ID:         utils.GenerateULID(),
		SenderType: models.SenderTypeAvatarAI,
		Message:    llmResponse.Message,
		CreatedAt:  time.Now(),
	}

	if err := avatarChatAdapter.CreateAvatarChatMessage(ctx, userID, avatarID, avatarResponse); err != nil {
		return nil, utils.WrapError(err)
	}

	newMatchingPoint := relation.MatchingPoint + llmResponse.PointChange
	if newMatchingPoint < 0 {
		newMatchingPoint = 0
	}
	if newMatchingPoint > 100 {
		newMatchingPoint = 100
	}

	relation.MatchingPoint = newMatchingPoint
	if err := avatarAdapter.UpdateUserAvatarRelation(relation); err != nil {
		return nil, utils.WrapError(err)
	}

	isMatched := isAlreadyMatched
	if !isAlreadyMatched && newMatchingPoint >= 100 {
		user1ID := userID
		user2ID := avatar.UserID
		if user1ID > user2ID {
			user1ID, user2ID = user2ID, user1ID
		}

		matching := models.Matching{
			ID:      utils.GenerateULID(),
			User1ID: user1ID,
			User2ID: user2ID,
		}

		if err := matchingAdapter.CreateMatching(matching); err != nil {
			log.Printf("Error creating matching: %v", err)
		} else {
			isMatched = true

			senderUser, err := userAdapter.GetByID(userID)
			if err != nil {
				log.Printf("Error getting sender user for notification: %v", err)
			} else {
				now := time.Now()

				notificationForSender := models.Notification{
					ID:        utils.GenerateULID(),
					UserID:    userID,
					Title:     "マッチング成立",
					Message:   fmt.Sprintf("%sさんとマッチングしました！", avatarOwnerUser.DisplayName),
					CreatedAt: now,
				}
				if err := notificationAdapter.CreateNotification(ctx, userID, notificationForSender); err != nil {
					log.Printf("Error creating notification for sender: %v", err)
				}

				notificationForOwner := models.Notification{
					ID:        utils.GenerateULID(),
					UserID:    avatar.UserID,
					Title:     "マッチング成立",
					Message:   fmt.Sprintf("%sさんとマッチングしました！", senderUser.DisplayName),
					CreatedAt: now,
				}
				if err := notificationAdapter.CreateNotification(ctx, avatar.UserID, notificationForOwner); err != nil {
					log.Printf("Error creating notification for avatar owner: %v", err)
				}
			}
		}
	}

	unlockedMissions, err := s.checkAndUnlockMissions(userID, avatar.UserID, newMatchingPoint, missionAdapter, userInfoAdapter)
	if err != nil {
		log.Printf("Error checking missions: %v", err)
	}

	return &SendMessageResult{
		AvatarResponse:   avatarResponse,
		MatchingPoint:    newMatchingPoint,
		PointChange:      llmResponse.PointChange,
		IsMatched:        isMatched,
		UnlockedMissions: unlockedMissions,
	}, nil
}

func (s *AvatarChatService) generateAvatarResponse(
	avatarOwnerUser models.User,
	avatarOwnerUserInfos []*models.UserInfo,
	chatHistory []adapter.AvatarChatMessage,
	llmAdapter adapter.LLMAdapter,
) (*LLMChatResponse, error) {
	userInfoStr := ""
	for _, info := range avatarOwnerUserInfos {
		userInfoStr += fmt.Sprintf("- %s: %s\n", info.Key, info.Value)
	}

	chatHistoryStr := ""
	for _, msg := range chatHistory {
		var sender string
		switch msg.SenderType {
		case models.SenderTypeUser:
			sender = "相手"
		case models.SenderTypeAvatarAI:
			sender = avatarOwnerUser.DisplayName
		default:
			sender = "システム"
		}
		chatHistoryStr += fmt.Sprintf("%s: 「%s」\n", sender, msg.Message)
	}

	prompt := fmt.Sprintf(`
# 命令
あなたは「%s」という名前のユーザーの分身AIです。
以下のユーザー情報を参考にして、相手からのメッセージに対して自然な会話を行ってください。

# ユーザー情報
名前: %s
性別: %s
自己紹介: %s

# 詳細情報
%s

# これまでの会話履歴
%s

# 出力形式
以下のJSON形式で出力してください。他の文字は一切出力しないでください。
- message: 相手へのメッセージ（2〜3文程度、日本語）
- point_change: 会話の質に基づくポイント変化（-10〜+10の整数）
  - 良い会話（共通点発見、質問への丁寧な回答、興味を示す）: +5〜+10
  - 普通の会話: +1〜+4
  - 微妙な会話（無関心、失礼な発言）: -5〜-10
- reason: ポイント変化の理由（簡潔に）

{
  "message": "メッセージ内容",
  "point_change": 5,
  "reason": "理由"
}
`, avatarOwnerUser.DisplayName, avatarOwnerUser.DisplayName, avatarOwnerUser.Gender, avatarOwnerUser.Bio, userInfoStr, chatHistoryStr)

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: prompt,
				},
			},
		},
	}

	resp, err := llmAdapter.CreateChatCompletion(contents, adapter.LLM_MODEL_TYPE_GEMINI2_5_FLASH)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var llmResponse LLMChatResponse
	if err := json.Unmarshal([]byte(resp), &llmResponse); err != nil {
		log.Printf("Failed to parse LLM response as JSON: %v, response: %s", err, resp)
		return &LLMChatResponse{
			Message:     resp,
			PointChange: 5,
			Reason:      "デフォルトポイント",
		}, nil
	}

	return &llmResponse, nil
}

func (s *AvatarChatService) checkAndUnlockMissions(
	userID string,
	avatarOwnerUserID string,
	currentMatchingPoint int,
	missionAdapter adapter.MissionAdapter,
	userInfoAdapter adapter.UserInfoAdapter,
) ([]UnlockedMissionInfo, error) {
	missions, err := missionAdapter.GetMissionsByOwnerUserID(avatarOwnerUserID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	var unlockedMissions []UnlockedMissionInfo

	for _, mission := range missions {
		if mission.ThresholdPointCondition == nil {
			continue
		}

		if currentMatchingPoint < *mission.ThresholdPointCondition {
			continue
		}

		existingUnlock, _ := missionAdapter.GetMissionUnlock(mission.ID, userID)
		if existingUnlock != nil {
			continue
		}

		missionUnlock := models.MissionUnlock{
			ID:             utils.GenerateULID(),
			MissionID:      mission.ID,
			UnlockedUserID: userID,
		}

		if err := missionAdapter.CreateMissionUnlock(missionUnlock); err != nil {
			log.Printf("Error creating mission unlock: %v", err)
			continue
		}

		userInfo, err := userInfoAdapter.GetByID(mission.UserInfoID)
		if err != nil {
			log.Printf("Error getting user info for mission: %v", err)
			continue
		}

		unlockedMissions = append(unlockedMissions, UnlockedMissionInfo{
			MissionID:  mission.ID,
			UserInfoID: mission.UserInfoID,
			Key:        userInfo.Key,
			Value:      userInfo.Value,
		})
	}

	return unlockedMissions, nil
}

func (s *AvatarChatService) GetMessages(ctx context.Context, userID string, avatarID string) ([]adapter.AvatarChatMessage, int, bool, error) {
	var avatarChatAdapter adapter.AvatarChatAdapter
	var avatarAdapter adapter.AvatarAdapter
	var matchingAdapter adapter.MatchingAdapter

	if err := s.container.Invoke(func(
		aca adapter.AvatarChatAdapter,
		aa adapter.AvatarAdapter,
		mta adapter.MatchingAdapter,
	) error {
		avatarChatAdapter = aca
		avatarAdapter = aa
		matchingAdapter = mta
		return nil
	}); err != nil {
		return nil, 0, false, utils.WrapError(err)
	}

	messages, err := avatarChatAdapter.GetAvatarChatMessages(ctx, userID, avatarID)
	if err != nil {
		return nil, 0, false, utils.WrapError(err)
	}

	avatar, err := avatarAdapter.GetByID(avatarID)
	if err != nil {
		return nil, 0, false, utils.WrapError(err)
	}

	relation, err := avatarAdapter.GetUserAvatarRelation(userID, avatarID)
	matchingPoint := 0
	if err == nil {
		matchingPoint = relation.MatchingPoint
	}

	existingMatching, _ := matchingAdapter.GetMatchingByUsers(userID, avatar.UserID)
	isMatched := existingMatching != nil

	return messages, matchingPoint, isMatched, nil
}

func (s *AvatarChatService) GetStatus(ctx context.Context, userID string, avatarID string) (int, bool, []UnlockedMissionInfo, error) {
	var avatarAdapter adapter.AvatarAdapter
	var matchingAdapter adapter.MatchingAdapter
	var missionAdapter adapter.MissionAdapter
	var userInfoAdapter adapter.UserInfoAdapter

	if err := s.container.Invoke(func(
		aa adapter.AvatarAdapter,
		mta adapter.MatchingAdapter,
		ma adapter.MissionAdapter,
		uia adapter.UserInfoAdapter,
	) error {
		avatarAdapter = aa
		matchingAdapter = mta
		missionAdapter = ma
		userInfoAdapter = uia
		return nil
	}); err != nil {
		return 0, false, nil, utils.WrapError(err)
	}

	avatar, err := avatarAdapter.GetByID(avatarID)
	if err != nil {
		return 0, false, nil, utils.WrapError(err)
	}

	relation, err := avatarAdapter.GetUserAvatarRelation(userID, avatarID)
	matchingPoint := 0
	if err == nil {
		matchingPoint = relation.MatchingPoint
	}

	existingMatching, _ := matchingAdapter.GetMatchingByUsers(userID, avatar.UserID)
	isMatched := existingMatching != nil

	missions, err := missionAdapter.GetMissionsByOwnerUserID(avatar.UserID)
	if err != nil {
		return matchingPoint, isMatched, nil, nil
	}

	var unlockedMissions []UnlockedMissionInfo
	for _, mission := range missions {
		existingUnlock, _ := missionAdapter.GetMissionUnlock(mission.ID, userID)
		if existingUnlock == nil {
			continue
		}

		userInfo, err := userInfoAdapter.GetByID(mission.UserInfoID)
		if err != nil {
			continue
		}

		unlockedMissions = append(unlockedMissions, UnlockedMissionInfo{
			MissionID:  mission.ID,
			UserInfoID: mission.UserInfoID,
			Key:        userInfo.Key,
			Value:      userInfo.Value,
		})
	}

	return matchingPoint, isMatched, unlockedMissions, nil
}
