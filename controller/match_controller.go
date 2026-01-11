package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type MatchController struct {
	container *dig.Container
}

func NewMatchController(container *dig.Container) *MatchController {
	return &MatchController{container: container}
}

// @Summary アンロック状態取得
// @Tags matches
// @Description 特定チャットのアンロック状態（マッチング可能かどうか）を取得する（IDは相手のUserID）
// @Security Bearer
// @Param partnerUserId path string true "相手のユーザーID（ULID）"
// @Success 200 {object} response.GetUnlockStatusResponse "アンロック状態取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このチャットにアクセスする権限がない"
// @Failure 404 {object} response.ErrorResponse "チャットが見つからない"
// @Router /chats/{partnerUserId}/unlock-status [get]
func (c *MatchController) GetUnlockStatus(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	partnerUserID := ctx.Param("partnerUserId")

	// TODO: 相手のユーザーIDからチャットが存在するか確認
	// TODO: ユーザーがこのチャットにアクセスする権限があるか確認
	// TODO: 既にマッチング済みでないか確認
	// TODO: 現在のマッチングスコアを取得（Firestoreから）
	// TODO: アンロック閾値と比較して状態を判定

	_ = userID        // TODO: 権限チェックで使用
	_ = partnerUserID // TODO: チャット情報取得で使用

	return ctx.JSON(http.StatusOK, &response.GetUnlockStatusResponse{
		Status: response.UnlockStatus{
			ChatID:          "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			IsUnlocked:      false,
			CurrentScore:    75,
			UnlockThreshold: 80,
			RemainingScore:  5,
			CanMatch:        false,
		},
	})
}

// @Summary マッチング成立（サーバー側で自動判定）
// @Tags matches
// @Description マッチング条件を達成した場合にマッチングを成立させる（サーバー側で自動判定、Firestore通知）
// @Security Bearer
// @Param partnerUserId path string true "相手のユーザーID（ULID）"
// @Success 200 {object} response.MatchResponse "マッチング成立成功"
// @Failure 400 {object} response.ErrorResponse "マッチング条件を満たしていない"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "このチャットにアクセスする権限がない、または既にマッチング済み"
// @Failure 404 {object} response.ErrorResponse "チャットが見つからない"
// @Router /chats/{partnerUserId}/match [post]
func (c *MatchController) Match(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)
	partnerUserID := ctx.Param("partnerUserId")

	// TODO: 相手のユーザーIDからチャットが存在するか確認
	// TODO: ユーザーがこのチャットにアクセスする権限があるか確認
	// TODO: 既にマッチング済みでないか確認
	// TODO: 現在のマッチングスコアを取得（Firestoreから）
	// TODO: マッチング条件（スコアが閾値以上）を満たしているか確認
	// TODO: マッチングを成立させる（マッチレコードを作成）
	// TODO: チャットのIsMatchedフラグを更新（Firestore）
	// TODO: Firestoreを通じて相手ユーザーに通知を送信

	return ctx.JSON(http.StatusOK, &response.MatchResponse{
		Match: response.Match{
			ID:              "01ARZ3NDEKTSV4RRFFQ69G5FEV",
			UserID:          userID,
			PartnerID:       partnerUserID,
			ChatID:          "01ARZ3NDEKTSV4RRFFQ69G5FAV",
			MatchedAt:       "2024-01-01T12:00:00Z",
			FinalScore:      85,
			PartnerName:     "佐藤花子",
			PartnerImageURL: "https://example.com/images/profile2.jpg",
		},
		Message: "マッチングが成立しました！",
	})
}
