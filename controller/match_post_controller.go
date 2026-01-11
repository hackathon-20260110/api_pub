package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type MatchPostController struct {
	container *dig.Container
}

func NewMatchPostController(container *dig.Container) *MatchPostController {
	return &MatchPostController{container: container}
}

// @Summary 相手候補一覧取得
// @Tags matches
// @Description マッチング前の相手候補の一覧を取得する（基本情報のみ）
// @Security Bearer
// @Param limit query int false "取得件数" default(20)
// @Param offset query int false "オフセット" default(0)
// @Success 200 {object} response.GetPartnersResponse "相手候補一覧取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /matches/candidates [get]
func (c *MatchPostController) GetCandidates(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)

	// TODO: ユーザーIDから相手候補一覧を取得
	// TODO: フィルタリング（年齢、性別など）

	_ = userID // TODO: 候補取得で使用

	return ctx.JSON(http.StatusOK, &response.GetPartnersResponse{
		Partners: []response.Partner{
			{
				ID:              "01ARZ3NDEKTSV4RRFFQ69G5FBV",
				DisplayName:     "佐藤花子",
				Age:             24,
				Gender:          "female",
				ProfileImageURL: "https://example.com/images/profile2.jpg",
				Bio:             "よろしくお願いします！",
			},
			{
				ID:              "01ARZ3NDEKTSV4RRFFQ69G5FCV",
				DisplayName:     "鈴木一郎",
				Age:             28,
				Gender:          "male",
				ProfileImageURL: "https://example.com/images/profile3.jpg",
				Bio:             "はじめまして！",
			},
		},
		Total:  50,
		Limit:  20,
		Offset: 0,
	})
}

// @Summary マッチ一覧取得
// @Tags matches
// @Description 自分がマッチングした相手の一覧を取得する
// @Security Bearer
// @Success 200 {object} response.GetMatchesResponse "マッチ一覧取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /matches [get]
func (c *MatchPostController) GetMatches(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)

	// TODO: ユーザーIDからマッチ一覧を取得
	// TODO: 各マッチの最新メッセージ、未読数などを取得

	_ = userID // TODO: マッチ一覧取得で使用

	return ctx.JSON(http.StatusOK, &response.GetMatchesResponse{
		Matches: []response.Match{
			{
				ID:              "01ARZ3NDEKTSV4RRFFQ69G5FEV",
				UserID:          userID,
				PartnerID:       "01ARZ3NDEKTSV4RRFFQ69G5FBV",
				ChatID:          "01ARZ3NDEKTSV4RRFFQ69G5FAV",
				MatchedAt:       "2024-01-01T12:00:00Z",
				FinalScore:      85,
				PartnerName:     "佐藤花子",
				PartnerImageURL: "https://example.com/images/profile2.jpg",
			},
		},
		Total: 1,
	})
}

// @Summary 話題提案取得
// @Tags matches
// @Description マッチした相手との会話を促進するための話題提案を取得する
// @Security Bearer
// @Param id path string true "マッチID"
// @Success 200 {object} response.GetSuggestionsResponse "話題提案取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このマッチにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "マッチが見つからない"
// @Router /matches/{id}/suggestions [get]
func (c *MatchPostController) GetSuggestions(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	matchID := ctx.Param("id")

	// TODO: マッチが存在するか確認
	// TODO: ユーザーがこのマッチにアクセスする権限があるか確認
	// TODO: ユーザーとパートナーのプロフィール情報を取得
	// TODO: 共通の興味、趣味、職業などを分析
	// TODO: AIを使って話題提案を生成

	_ = userID  // TODO: 権限チェックで使用
	_ = matchID // TODO: マッチ情報取得で使用

	return ctx.JSON(http.StatusOK, &response.GetSuggestionsResponse{
		Suggestions: []response.Suggestion{
			{
				ID:          "suggestion_123456",
				Title:       "共通の趣味について話す",
				Description: "お互いの趣味について詳しく聞いてみましょう",
				Category:    "interests",
				Examples:    []string{"どんな本が好きですか？", "最近読んだ本はありますか？"},
			},
			{
				ID:          "suggestion_123457",
				Title:       "仕事について話す",
				Description: "お互いの仕事について話してみましょう",
				Category:    "work",
				Examples:    []string{"どんなお仕事をされていますか？", "仕事の楽しさは何ですか？"},
			},
		},
	})
}

// @Summary 本人とのメッセージ送信
// @Tags matches
// @Description マッチした相手（本人）にメッセージを送信する
// @Security Bearer
// @Param id path string true "マッチID"
// @Param request body requests.SendMatchMessageRequest true "メッセージ送信リクエスト"
// @Success 200 {object} response.SendMatchMessageResponse "メッセージ送信成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このマッチにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "マッチが見つからない"
// @Router /matches/{id}/messages [post]
func (c *MatchPostController) SendMatchMessage(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	matchID := ctx.Param("id")

	// TODO: リクエストボディをパース
	// var req requests.SendMatchMessageRequest
	// if err := ctx.Bind(&req); err != nil {
	//     return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
	//         Error:   "invalid_request",
	//         Message: "リクエストが不正です",
	//     })
	// }

	// TODO: マッチが存在するか確認
	// TODO: ユーザーがこのマッチにアクセスする権限があるか確認
	// TODO: メッセージを保存
	// TODO: 相手ユーザーに通知を送信（必要に応じて）

	_ = userID  // TODO: 権限チェックで使用
	_ = matchID // TODO: マッチ情報取得で使用

	return ctx.JSON(http.StatusOK, &response.SendMatchMessageResponse{
		Message: response.MatchMessage{
			ID:        "match_msg_123456",
			MatchID:   matchID,
			SenderID:  userID,
			Content:   "こんにちは！マッチできて嬉しいです。",
			CreatedAt: "2024-01-01T12:00:00Z",
		},
	})
}

// @Summary 本人とのメッセージ履歴取得
// @Tags matches
// @Description マッチした相手（本人）とのメッセージ履歴を取得する
// @Security Bearer
// @Param id path string true "マッチID"
// @Param limit query int false "取得件数" default(20)
// @Param offset query int false "オフセット" default(0)
// @Success 200 {object} response.GetMatchMessagesResponse "メッセージ履歴取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このマッチにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "マッチが見つからない"
// @Router /matches/{id}/messages [get]
func (c *MatchPostController) GetMatchMessages(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	matchID := ctx.Param("id")

	// TODO: マッチが存在するか確認
	// TODO: ユーザーがこのマッチにアクセスする権限があるか確認
	// TODO: メッセージ履歴を取得（ページネーション対応）

	_ = userID  // TODO: 権限チェックで使用
	_ = matchID // TODO: マッチ情報取得で使用

	return ctx.JSON(http.StatusOK, &response.GetMatchMessagesResponse{
		Messages: []response.MatchMessage{
			{
				ID:        "match_msg_123456",
				MatchID:   matchID,
				SenderID:  userID,
				Content:   "こんにちは！マッチできて嬉しいです。",
				CreatedAt: "2024-01-01T12:00:00Z",
			},
			{
				ID:        "match_msg_123457",
				MatchID:   matchID,
				SenderID:  "01ARZ3NDEKTSV4RRFFQ69G5FBV",
				Content:   "こちらこそ！よろしくお願いします。",
				CreatedAt: "2024-01-01T12:05:00Z",
			},
		},
		Total:  2,
		Limit:  20,
		Offset: 0,
	})
}

// @Summary 返信アシスト
// @Tags matches
// @Description 相手からの最新メッセージに対する返信をAIが提案する（チャットルーム入室時や返信受信時に自動取得）
// @Security Bearer
// @Param id path string true "マッチID（ULID）"
// @Param messageId query string false "返信対象のメッセージID（指定しない場合は最新メッセージ）"
// @Success 200 {object} response.ReplyAssistResponse "返信アシスト取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このマッチにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "マッチが見つからない"
// @Router /matches/{id}/reply-assist [get]
func (c *MatchPostController) ReplyAssist(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	matchID := ctx.Param("id")
	messageID := ctx.QueryParam("messageId") // オプション: 特定のメッセージに対する返信

	// TODO: マッチが存在するか確認
	// TODO: ユーザーがこのマッチにアクセスする権限があるか確認
	// TODO: 会話履歴を取得（Firestoreから）
	// TODO: 指定されたメッセージIDまたは最新メッセージを取得
	// TODO: AIを使って返信候補を生成（相手のメッセージ、会話の文脈、ユーザーのプロフィールなどを考慮）

	_ = userID    // TODO: 権限チェックで使用
	_ = matchID   // TODO: マッチ情報取得で使用
	_ = messageID // TODO: メッセージ取得で使用

	return ctx.JSON(http.StatusOK, &response.ReplyAssistResponse{
		SuggestedReplies: []string{
			"ありがとうございます！こちらこそよろしくお願いします。",
			"こちらこそ！お話できるのを楽しみにしています。",
			"ありがとうございます！早速お話しできて嬉しいです。",
		},
		Context: "初めてのメッセージに対する返信",
	})
}
