package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/hackathon-20260110/api/models"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/utils"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

type DiagnosisService interface {
	// Avatar同士の会話を生成して診断を実行
	ExecuteDiagnosis(userID, targetAvatarID string, conversationData map[string]interface{}) (*response.DiagnosisResult, error)

	// 診断履歴の取得
	GetDiagnosisHistory(userID string) ([]response.DiagnosisHistory, error)

	// 特定の診断結果取得
	GetDiagnosisDetail(userID, diagnosisID string) (*response.DiagnosisDetail, error)
}

type diagnosisService struct {
	diagnosisAdapter adapter.DiagnosisAdapter
	llmAdapter       adapter.LLMAdapter
}

func NewDiagnosisService(diagnosisAdapter adapter.DiagnosisAdapter, llmAdapter adapter.LLMAdapter) DiagnosisService {
	return &diagnosisService{
		diagnosisAdapter: diagnosisAdapter,
		llmAdapter:       llmAdapter,
	}
}

// ポイント計算マップ
var scoreToPoints = map[int]int{
	1: 0,
	2: 20,
	3: 40,
	4: 60,
	5: 100,
}

func (s *diagnosisService) ExecuteDiagnosis(userID, targetAvatarID string, conversationData map[string]interface{}) (*response.DiagnosisResult, error) {
	// Avatarの存在確認
	userAvatar, err := s.diagnosisAdapter.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("user avatar not found: %w", err)
	}

	targetAvatar, err := s.diagnosisAdapter.GetAvatarByID(targetAvatarID)
	if err != nil {
		return nil, fmt.Errorf("target avatar not found: %w", err)
	}

	// 会話データをJSON形式に変換
	conversationJSON, err := json.Marshal(conversationData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal conversation data: %w", err)
	}

	// AI診断を実行
	diagnosisScore, analysisResult, err := s.performAIDiagnosis(userAvatar, targetAvatar, conversationData)
	if err != nil {
		return nil, fmt.Errorf("AI diagnosis failed: %w", err)
	}

	// 分析結果をJSON形式に変換
	analysisJSON, err := json.Marshal(analysisResult)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analysis result: %w", err)
	}

	// 診断履歴を保存
	diagnosisHistory := models.DiagnosisHistory{
		ID:               utils.GenerateULID(),
		UserID:           userID,
		UserAvatarID:     userAvatar.ID,
		TargetAvatarID:   targetAvatarID,
		ConversationData: string(conversationJSON),
		DiagnosisScore:   diagnosisScore,
		AIAnalysisResult: string(analysisJSON),
	}

	if err := s.diagnosisAdapter.CreateDiagnosisHistory(&diagnosisHistory); err != nil {
		return nil, fmt.Errorf("failed to save diagnosis history: %w", err)
	}

	// ポイント加算処理
	pointsEarned := scoreToPoints[diagnosisScore]
	if err := s.updateMatchingPoints(userID, targetAvatarID, pointsEarned); err != nil {
		return nil, fmt.Errorf("failed to update matching points: %w", err)
	}

	// レスポンス作成
	return &response.DiagnosisResult{
		DiagnosisID:    diagnosisHistory.ID,
		DiagnosisScore: diagnosisScore,
		PointsEarned:   pointsEarned,
		AnalysisResult: analysisResult,
		CanDirectChat:  diagnosisScore == 5, // スコア5で直接チャット解禁
		CreatedAt:      diagnosisHistory.CreatedAt,
	}, nil
}

func (s *diagnosisService) GetDiagnosisHistory(userID string) ([]response.DiagnosisHistory, error) {
	histories, err := s.diagnosisAdapter.GetDiagnosisHistoryByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnosis history: %w", err)
	}

	var result []response.DiagnosisHistory
	for _, history := range histories {
		var analysisResult map[string]interface{}
		if err := json.Unmarshal([]byte(history.AIAnalysisResult), &analysisResult); err != nil {
			analysisResult = map[string]interface{}{"error": "failed to parse analysis result"}
		}

		result = append(result, response.DiagnosisHistory{
			ID:               history.ID,
			UserAvatarName:   history.UserAvatar.Prompt, // Avatarの名前として使用
			TargetAvatarName: history.TargetAvatar.Prompt,
			DiagnosisScore:   history.DiagnosisScore,
			PointsEarned:     scoreToPoints[history.DiagnosisScore],
			CanDirectChat:    history.DiagnosisScore == 5,
			CreatedAt:        history.CreatedAt,
		})
	}

	return result, nil
}

func (s *diagnosisService) GetDiagnosisDetail(userID, diagnosisID string) (*response.DiagnosisDetail, error) {
	history, err := s.diagnosisAdapter.GetDiagnosisHistoryByID(diagnosisID)
	if err != nil {
		return nil, fmt.Errorf("diagnosis not found: %w", err)
	}

	// 権限チェック
	if history.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	var conversationData map[string]interface{}
	if err := json.Unmarshal([]byte(history.ConversationData), &conversationData); err != nil {
		return nil, fmt.Errorf("failed to parse conversation data: %w", err)
	}

	var analysisResult map[string]interface{}
	if err := json.Unmarshal([]byte(history.AIAnalysisResult), &analysisResult); err != nil {
		analysisResult = map[string]interface{}{"error": "failed to parse analysis result"}
	}

	return &response.DiagnosisDetail{
		ID:               history.ID,
		UserAvatar:       response.NewAvatarResponse(history.UserAvatar),
		TargetAvatar:     response.NewAvatarResponse(history.TargetAvatar),
		ConversationData: conversationData,
		DiagnosisScore:   history.DiagnosisScore,
		PointsEarned:     scoreToPoints[history.DiagnosisScore],
		AnalysisResult:   analysisResult,
		CanDirectChat:    history.DiagnosisScore == 5,
		CreatedAt:        history.CreatedAt,
	}, nil
}

// sanitizeJSONResponse はLLMレスポンスからJSONを抽出・正規化する
func sanitizeJSONResponse(resp string) string {
	resp = strings.TrimSpace(resp)
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	// まだJSONとして有効でない場合、先頭の{から末尾の}までを抽出
	startIdx := strings.Index(resp, "{")
	endIdx := strings.LastIndex(resp, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		resp = resp[startIdx : endIdx+1]
	}

	return resp
}

// AI診断を実行（LLMAdapterを使用）
func (s *diagnosisService) performAIDiagnosis(userAvatar, targetAvatar models.Avatar, conversationData map[string]interface{}) (int, map[string]interface{}, error) {
	// 会話データをJSON文字列に変換
	conversationJSON, err := json.Marshal(conversationData)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to marshal conversation data: %w", err)
	}

	// 診断プロンプトを構築
	prompt := fmt.Sprintf(`
以下のアバター同士の会話を分析して、相性を1-5段階で評価してください。

【アバター1】
プロンプト: %s
性格特性: %s

【アバター2】
プロンプト: %s
性格特性: %s

【会話内容】
%s

相性スコア（1-5）と理由を以下のJSON形式のみで回答してください。JSON以外の文字は出力しないでください：
{
  "score": 数値,
  "reason": "詳細な理由",
  "compatibility_factors": ["要因1", "要因2", ...],
  "improvement_suggestions": ["提案1", "提案2", ...]
}
`, userAvatar.Prompt, userAvatar.PersonalityTraits, targetAvatar.Prompt, targetAvatar.PersonalityTraits, string(conversationJSON))

	// LLMで分析実行（JSON出力指定）
	contents := genai.Text(prompt)
	aiResponse, err := s.llmAdapter.CreateChatCompletionJSON(contents, adapter.LLM_MODEL_TYPE_GEMINI2_0)
	if err != nil {
		return 0, nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	// AIレスポンスをサニタイズしてパース
	sanitized := sanitizeJSONResponse(aiResponse)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(sanitized), &result); err != nil {
		// JSONパースに失敗した場合のフォールバック
		return 3, map[string]interface{}{
			"score":  3,
			"reason": "AI分析でエラーが発生しました",
			"error":  err.Error(),
		}, nil
	}

	score, ok := result["score"].(float64)
	if !ok || score < 1 || score > 5 {
		score = 3 // デフォルトスコア
	}

	return int(score), result, nil
}

// マッチングポイントを更新
func (s *diagnosisService) updateMatchingPoints(userID, avatarID string, pointsToAdd int) error {
	if pointsToAdd == 0 {
		return nil // ポイント加算なしの場合はスキップ
	}

	// 既存のUserAvatarRelationを確認
	_, err := s.diagnosisAdapter.GetUserAvatarRelation(userID, avatarID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 新規作成
			relation := models.UserAvatarRelation{
				ID:            utils.GenerateULID(),
				UserID:        userID,
				AvatarID:      avatarID,
				MatchingPoint: pointsToAdd,
			}
			return s.diagnosisAdapter.CreateUserAvatarRelation(relation)
		}
		return err
	}

	// 既存レコードを更新
	return s.diagnosisAdapter.UpdateUserAvatarRelationPoints(userID, avatarID, pointsToAdd)
}
