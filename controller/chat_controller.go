package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type ChatController struct {
	container *dig.Container
}

func NewChatController(container *dig.Container) *ChatController {
	return &ChatController{container: container}
}

// @Summary 新しいチャット開始
// @Tags chats
// @Description 相手との新しいチャットを開始する
// @Security Bearer
// @Param request body requests.CreateChatRequest true "チャット開始リクエスト"
// @Success 200 {object} response.CreateChatResponse "チャット開始成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正、または既にチャットが存在する"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 404 {object} response.ErrorResponse "相手が見つからない"
// @Router /chats [post]
func (c *ChatController) CreateChat(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)

	// TODO: リクエストボディをパース
	// var req requests.CreateChatRequest
	// if err := ctx.Bind(&req); err != nil {
	//     return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
	//         Error:   "invalid_request",
	//         Message: "リクエストが不正です",
	//     })
	// }

	// TODO: 相手が存在するか確認
	// TODO: 既にチャットが存在しないか確認
	// TODO: チャットを作成
	// TODO: 初期メッセージがある場合は送信

	return ctx.JSON(http.StatusOK, &response.CreateChatResponse{
		Chat: response.Chat{
			ID:              "chat_123456",
			UserID:          userID,
			PartnerID:       "01ARZ3NDEKTSV4RRFFQ69G5FBV",
			PartnerName:     "佐藤花子",
			PartnerImageURL: "https://example.com/images/profile2.jpg",
			MatchingScore:   0,
			IsMatched:       false,
			CreatedAt:       "2024-01-01T00:00:00Z",
			UpdatedAt:       "2024-01-01T00:00:00Z",
		},
	})
}

// @Summary 自分のチャット一覧取得
// @Tags chats
// @Description 自分が参加しているチャットの一覧を取得する
// @Security Bearer
// @Success 200 {object} response.GetChatsResponse "チャット一覧取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /chats [get]
func (c *ChatController) GetChats(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)

	chatService := service.NewChatService(c.container)
	result, err := chatService.GetChats(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_error",
			Message: "チャット一覧の取得に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary 特定チャットの詳細取得
// @Tags chats
// @Description 特定のチャットの詳細情報を取得する（IDは相手のUserID）
// @Security Bearer
// @Param partnerUserId path string true "相手のユーザーID（ULID）"
// @Success 200 {object} response.GetChatDetailResponse "チャット詳細取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このチャットにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "チャットが見つからない"
// @Router /chats/{partnerUserId} [get]
func (c *ChatController) GetChatDetail(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	partnerUserID := ctx.Param("partnerUserId")

	// TODO: 相手のユーザーIDからチャットが存在するか確認
	// TODO: ユーザーがこのチャットにアクセスする権限があるか確認
	// TODO: チャット詳細を取得（Firestoreから）

	return ctx.JSON(http.StatusOK, &response.GetChatDetailResponse{
		Chat: response.Chat{
			ID:              "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			UserID:          userID,
			PartnerID:       partnerUserID,
			PartnerName:     "佐藤花子",
			PartnerImageURL: "https://example.com/images/profile2.jpg",
			LastMessage:     "こんにちは！",
			LastMessageAt:   "2024-01-01T12:00:00Z",
			UnreadCount:     0,
			MatchingScore:   75,
			IsMatched:       false,
			CreatedAt:       "2024-01-01T00:00:00Z",
			UpdatedAt:       "2024-01-01T12:00:00Z",
		},
	})
}

// @Summary メッセージ送信
// @Tags chats
// @Description 相手AIにメッセージを送信する（IDは相手のUserID、メッセージ履歴はFirestoreから直接取得）
// @Security Bearer
// @Param partnerUserId path string true "相手のユーザーID（ULID）"
// @Param request body requests.SendMessageRequest true "メッセージ送信リクエスト"
// @Success 200 {object} response.SendMessageResponse "メッセージ送信成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このチャットにアクセスする権限がない、またはマッチング済み"
// @Failure 404 {object} response.ErrorResponse "チャットが見つからない"
// @Router /chats/{partnerUserId}/messages [post]
func (c *ChatController) SendMessage(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	partnerUserID := ctx.Param("partnerUserId")

	// TODO: リクエストボディをパース
	// var req requests.SendMessageRequest
	// if err := ctx.Bind(&req); err != nil {
	//     return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
	//         Error:   "invalid_request",
	//         Message: "リクエストが不正です",
	//     })
	// }

	// TODO: 相手のユーザーIDからチャットが存在するか確認
	// TODO: ユーザーがこのチャットにアクセスする権限があるか確認
	// TODO: マッチング済みでないか確認（マッチング後は別のエンドポイントを使用）
	// TODO: メッセージをFirestoreに保存
	// TODO: 相手の分身AIにメッセージを送信してAI返信を取得
	// TODO: AI返信をFirestoreに保存
	// TODO: マッチングスコアを更新（会話の質、共通の興味、エンゲージメントなど）
	// TODO: マッチング条件を満たした場合は自動的にマッチング成立処理を実行（Firestore通知）

	_ = partnerUserID // TODO: チャット情報取得で使用

	return ctx.JSON(http.StatusOK, &response.SendMessageResponse{
		Message: response.Message{
			ID:         "01ARZ3NDEKTSV4RRFFQ69G5FBV",
			ChatID:     "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			SenderID:   userID,
			SenderType: "user",
			Content:    "今日はいい天気ですね。",
			CreatedAt:  "2024-01-01T12:00:10Z",
		},
		AIResponse: &response.Message{
			ID:         "01ARZ3NDEKTSV4RRFFQ69G5FCV",
			ChatID:     "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			SenderID:   "01ARZ3NDEKTSV4RRFFQ69G5FDV",
			SenderType: "avatar_ai",
			Content:    "そうですね！こんな日はお散歩が気持ちいいですよね。",
			CreatedAt:  "2024-01-01T12:00:15Z",
		},
		MatchingScore: 75,
	})
}

// @Summary 現在のマッチングポイント取得
// @Tags chats
// @Description 特定チャットの現在のマッチングポイントを取得する（IDは相手のUserID）
// @Security Bearer
// @Param partnerUserId path string true "相手のユーザーID（ULID）"
// @Success 200 {object} response.GetChatScoreResponse "マッチングポイント取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このチャットにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "チャットが見つからない"
// @Router /chats/{partnerUserId}/score [get]
func (c *ChatController) GetChatScore(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	partnerUserID := ctx.Param("partnerUserId")

	// TODO: 相手のユーザーIDからチャットが存在するか確認
	// TODO: ユーザーがこのチャットにアクセスする権限があるか確認
	// TODO: マッチングスコアを取得（Firestoreから）

	_ = userID        // TODO: 権限チェックで使用
	_ = partnerUserID // TODO: チャット情報取得で使用

	return ctx.JSON(http.StatusOK, &response.GetChatScoreResponse{
		Score: response.ChatScore{
			ChatID:       "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			CurrentScore: 75,
			MaxScore:     100,
			ScoreBreakdown: map[string]int{
				"conversation_quality": 30,
				"common_interests":     25,
				"engagement":           20,
			},
			UnlockThreshold: 80,
			UpdatedAt:       "2024-01-01T12:00:00Z",
		},
	})
}
