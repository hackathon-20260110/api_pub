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
)

type OnboardingService struct {
	container *dig.Container
}

func NewOnboardingService(container *dig.Container) *OnboardingService {
	return &OnboardingService{container: container}
}

func (s *OnboardingService) StartOnboardingChat(ctx context.Context, userID string) error {
	var onboardingAdapter adapter.OnboardingAdapter
	var llmAdapter adapter.LLMAdapter
	var userAdapter adapter.UserAdapter
	if err := s.container.Invoke(func(oa adapter.OnboardingAdapter, la adapter.LLMAdapter, ua adapter.UserAdapter) error {
		onboardingAdapter = oa
		llmAdapter = la
		userAdapter = ua
		return nil
	}); err != nil {
		return utils.WrapError(err)
	}

	u, err := userAdapter.GetByID(userID)
	if err != nil {
		return utils.WrapError(err)
	}

	systemPrompt := `
# 命令
 - マッチングアプリにおける登録ユーザの情報をチャットで取得するための対話相手になってもらいます。
 - ユーザの趣味や特技、恋愛観、好きなタイプ、過去の恋愛歴などをチャットを通じて掘り下げます。
 - そのための最初のメッセージ書いてください。
 - メッセージは2~3文程度で作成してください。
 - メッセージは日本語で作成してください。
`
	userPrompt := fmt.Sprintf(
		`
# ユーザーの情報
次の情報がユーザのメタ情報である
名前: %s
性別: %s
生年月日: %s
自己紹介: %s`,
		u.DisplayName, u.Gender, u.BirthDate.Format("2006-01-02"), u.Bio,
	)

	content := []*genai.Content{
		{
			Role: "model",
			Parts: []*genai.Part{
				{
					Text: systemPrompt,
				},
			},
		},
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: userPrompt,
				},
			},
		},
	}

	resp, err := llmAdapter.CreateChatCompletion(content, adapter.LLM_MODEL_TYPE_GEMINI2_5_FLASH)
	if err != nil {
		return utils.WrapError(err)
	}

	chat := models.OnboardingChat{
		ID:         utils.GenerateULID(),
		SenderType: models.SenderTypeSystem,
		Message:    resp,
		CreatedAt:  time.Now(),
	}

	return onboardingAdapter.CreateOnboardingChat(ctx, userID, chat)
}

func (s *OnboardingService) SendOnboardingMessage(ctx context.Context, userID string, userMessage string) error {
	var onboardingAdapter adapter.OnboardingAdapter
	var userAdapter adapter.UserAdapter
	var llmAdapter adapter.LLMAdapter
	if err := s.container.Invoke(func(oa adapter.OnboardingAdapter, ua adapter.UserAdapter, la adapter.LLMAdapter) error {
		onboardingAdapter = oa
		userAdapter = ua
		llmAdapter = la
		return nil
	}); err != nil {
		return utils.WrapError(err)
	}

	u, err := userAdapter.GetByID(userID)
	if err != nil {
		return utils.WrapError(err)
	}

	chats, err := onboardingAdapter.GetOnboardingChats(ctx, userID)
	if err != nil {
		return utils.WrapError(err)
	}

	userChat := models.OnboardingChat{
		ID:         utils.GenerateULID(),
		SenderType: models.SenderTypeUser,
		Message:    userMessage,
		CreatedAt:  time.Now(),
	}

	if err := onboardingAdapter.CreateOnboardingChat(ctx, userID, userChat); err != nil {
		return utils.WrapError(err)
	}

	chats = append(chats, userChat)

	go s.processOnboardingMessageAsync(userID, &u, chats, onboardingAdapter, llmAdapter)

	return nil
}

func (s *OnboardingService) processOnboardingMessageAsync(
	userID string,
	u *models.User,
	chats []models.OnboardingChat,
	onboardingAdapter adapter.OnboardingAdapter,
	llmAdapter adapter.LLMAdapter,
) {
	ctx := context.Background()

	userMessageCount := 0
	for _, chat := range chats {
		if chat.SenderType == models.SenderTypeUser {
			userMessageCount++
		}
	}
	rallyCount := userMessageCount

	isOnboardingCompleted := false

	if rallyCount > 10 {
		isOnboardingCompleted = true
	} else {
		isOnboardingCompleted = s.judgeOnboardingCompletion(chats, llmAdapter)
	}

	systemPrompt := `
# 命令
 - マッチングアプリにおける登録ユーザの情報をチャットで取得するための対話相手になってもらいます。
 - ユーザの趣味や特技、恋愛観、好きなタイプ、過去の恋愛歴などをチャットを通じて掘り下げます。
 - これまでのメッセージ履歴を参考にして次のメッセージを作成してください。
 - メッセージは2~3文程度で作成してください。
 - メッセージは日本語で作成してください。
`

	chatHists := "# メッセージ履歴\n"
	for _, chat := range chats {
		var userName string
		switch chat.SenderType {
		case models.SenderTypeSystem:
			userName = "システム"
		case models.SenderTypeUser:
			userName = u.DisplayName
		}
		chatHists += fmt.Sprintf("%s: 「%s」\n", userName, chat.Message)
	}

	fullPrompt := systemPrompt + "\n" + chatHists

	content := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: fullPrompt,
				},
			},
		},
	}

	resp, err := llmAdapter.CreateChatCompletion(content, adapter.LLM_MODEL_TYPE_GEMINI2_5_FLASH)
	if err != nil {
		log.Printf("Error generating LLM response for user %s: %v", userID, err)
		return
	}

	systemChat := models.OnboardingChat{
		ID:                    utils.GenerateULID(),
		SenderType:            models.SenderTypeSystem,
		Message:               resp,
		IsOnboardingCompleted: isOnboardingCompleted,
		CreatedAt:             time.Now(),
	}

	if err := onboardingAdapter.CreateOnboardingChat(ctx, userID, systemChat); err != nil {
		log.Printf("Error saving system chat for user %s: %v", userID, err)
		return
	}
}

func (s *OnboardingService) judgeOnboardingCompletion(chats []models.OnboardingChat, llmAdapter adapter.LLMAdapter) bool {
	chatHistory := ""
	for _, chat := range chats {
		var sender string
		switch chat.SenderType {
		case models.SenderTypeSystem:
			sender = "システム"
		case models.SenderTypeUser:
			sender = "ユーザー"
		}
		chatHistory += fmt.Sprintf("%s: 「%s」\n", sender, chat.Message)
	}

	judgePrompt := fmt.Sprintf(`
# 命令
あなたはマッチングアプリのオンボーディングプロセスを評価するジャッジです。
以下のチャット履歴を分析し、ユーザーの分身AI（アバター）を作成するのに十分な情報が集まったかどうかを判断してください。

# 判断基準
アバターを作成するには、以下の情報がある程度揃っている必要があります：
- 趣味や興味関心
- 性格や価値観
- 恋愛観や好きなタイプ
- ライフスタイル（仕事、休日の過ごし方など）

# チャット履歴
%s

# 出力形式
十分な情報が集まっている場合は「TRUE」、まだ不十分な場合は「FALSE」とだけ出力してください。
他の文字は一切出力しないでください。
`, chatHistory)

	content := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: judgePrompt,
				},
			},
		},
	}

	resp, err := llmAdapter.CreateChatCompletion(content, adapter.LLM_MODEL_TYPE_GEMINI2_5_FLASH)
	if err != nil {
		log.Printf("Error in LLM-as-a-judge: %v", err)
		return false
	}

	result := strings.TrimSpace(strings.ToUpper(resp))
	return result == "TRUE"
}

func (s *OnboardingService) FinishOnboarding(ctx context.Context, userID string) ([]*models.UserInfo, *models.Avatar, error) {
	var onboardingAdapter adapter.OnboardingAdapter
	var userAdapter adapter.UserAdapter
	var avatarAdapter adapter.AvatarAdapter
	var userInfoAdapter adapter.UserInfoAdapter
	var llmAdapter adapter.LLMAdapter
	if err := s.container.Invoke(func(oa adapter.OnboardingAdapter, ua adapter.UserAdapter, aa adapter.AvatarAdapter, uia adapter.UserInfoAdapter, la adapter.LLMAdapter) error {
		onboardingAdapter = oa
		userAdapter = ua
		avatarAdapter = aa
		userInfoAdapter = uia
		llmAdapter = la
		return nil
	}); err != nil {
		return nil, nil, utils.WrapError(err)
	}

	u, err := userAdapter.GetByID(userID)
	if err != nil {
		return nil, nil, utils.WrapError(err)
	}

	chats, err := onboardingAdapter.GetOnboardingChats(ctx, userID)
	if err != nil {
		return nil, nil, utils.WrapError(err)
	}

	chatHistory := ""
	for _, chat := range chats {
		var sender string
		switch chat.SenderType {
		case models.SenderTypeSystem:
			sender = "システム"
		case models.SenderTypeUser:
			sender = "ユーザー"
		}
		chatHistory += fmt.Sprintf("%s: 「%s」\n", sender, chat.Message)
	}

	userInfos, err := s.extractUserInfoFromChat(chatHistory, userID, llmAdapter)
	if err != nil {
		return nil, nil, utils.WrapError(err)
	}

	if len(userInfos) > 0 {
		_, err = userInfoAdapter.CreateMany(userInfos)
		if err != nil {
			return nil, nil, utils.WrapError(err)
		}
	}

	avatar := models.Avatar{
		ID:                utils.GenerateULID(),
		UserID:            userID,
		AvatarIconURL:     u.ProfileImageURL,
		Prompt:            "",
		PersonalityTraits: "{}",
	}

	_, err = avatarAdapter.Create(avatar)
	if err != nil {
		return nil, nil, utils.WrapError(err)
	}

	u.IsOnboardingCompleted = true
	_, err = userAdapter.Update(u)
	if err != nil {
		return nil, nil, utils.WrapError(err)
	}

	return userInfos, &avatar, nil
}

func (s *OnboardingService) extractUserInfoFromChat(chatHistory string, userID string, llmAdapter adapter.LLMAdapter) ([]*models.UserInfo, error) {
	extractPrompt := fmt.Sprintf(`
# 命令
あなたはマッチングアプリのオンボーディングチャット履歴を分析し、ユーザーの情報を抽出するアシスタントです。
以下のチャット履歴から、ユーザーのプロフィール情報を抽出してください。

# 抽出する情報
- 趣味・興味関心
- 性格・価値観
- 恋愛観・好きなタイプ
- ライフスタイル（仕事、休日の過ごし方など）
- その他、マッチングに役立つ情報

# チャット履歴
%s

# 出力形式
以下のJSON形式で出力してください。各項目はkey-valueのペアで、keyは項目名、valueはユーザーが回答した内容です。
必ず有効なJSONのみを出力し、他の文字は一切出力しないでください。

{
  "items": [
    {"key": "趣味", "value": "映画鑑賞、読書"},
    {"key": "性格", "value": "明るく社交的"},
    {"key": "恋愛観", "value": "誠実な関係を大切にしたい"}
  ]
}
`, chatHistory)

	content := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: extractPrompt,
				},
			},
		},
	}

	resp, err := llmAdapter.CreateChatCompletion(content, adapter.LLM_MODEL_TYPE_GEMINI2_5_FLASH)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var result struct {
		Items []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		log.Printf("Failed to parse LLM response as JSON: %v, response: %s", err, resp)
		return nil, nil
	}

	userInfos := make([]*models.UserInfo, 0, len(result.Items))
	for _, item := range result.Items {
		userInfo := &models.UserInfo{
			ID:       utils.GenerateULID(),
			UserID:   userID,
			InfoType: models.UserInfoTypeText,
			Key:      item.Key,
			Value:    item.Value,
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}
