package controller

import (
	"fmt"
	"net/http"

	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
)

type DiagnosisController struct {
	diagnosisService service.DiagnosisService
}

func NewDiagnosisController(diagnosisService service.DiagnosisService) *DiagnosisController {
	return &DiagnosisController{
		diagnosisService: diagnosisService,
	}
}

// @Summary Avatar間診断実行
// @Tags diagnosis
// @Description 自分のAvatarと相手のAvatarの会話を分析して相性診断を実行し、マッチングポイントを加算する
// @Security Bearer
// @Param body body requests.ExecuteDiagnosisRequest true "診断リクエスト"
// @Success 200 {object} response.DiagnosisResult "診断実行成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 404 {object} response.ErrorResponse "Avatarが見つからない"
// @Failure 500 {object} response.ErrorResponse "サーバーエラー"
// @Router /diagnosis/execute [post]
func (c *DiagnosisController) ExecuteDiagnosis(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}
	var req requests.ExecuteDiagnosisRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Message: "リクエストが不正です",
		})
	}

	fmt.Println("req", req)

	result, err := c.diagnosisService.ExecuteDiagnosis(
		userID,
		req.TargetAvatarID,
		req.ConversationData,
	)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Message: "診断に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary 診断履歴取得
// @Tags diagnosis
// @Description ユーザーの診断履歴一覧を取得する
// @Security Bearer
// @Param limit query int false "取得件数" default(20)
// @Param offset query int false "オフセット" default(0)
// @Success 200 {object} response.DiagnosisHistoryResponse "診断履歴取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 500 {object} response.ErrorResponse "サーバーエラー"
// @Router /diagnosis/history [get]
func (c *DiagnosisController) GetDiagnosisHistory(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)
	if userID == "" {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	histories, err := c.diagnosisService.GetDiagnosisHistory(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Message: "取得に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, &response.DiagnosisHistoryResponse{
		Histories: histories,
	})
}

// @Summary 診断詳細取得
// @Tags diagnosis
// @Description 特定の診断結果の詳細情報を取得する
// @Security Bearer
// @Param diagnosisId path string true "診断ID（ULID）"
// @Success 200 {object} response.DiagnosisDetailResponse "診断詳細取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "アクセス権限がない"
// @Failure 404 {object} response.ErrorResponse "診断結果が見つからない"
// @Failure 500 {object} response.ErrorResponse "サーバーエラー"
// @Router /diagnosis/{diagnosisId} [get]
func (c *DiagnosisController) GetDiagnosisDetail(ctx echo.Context) error {
	userID := middleware.GetFirebaseUID(ctx)
	if userID == "" {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}
	diagnosisID := ctx.Param("diagnosisId")

	detail, err := c.diagnosisService.GetDiagnosisDetail(userID, diagnosisID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Message: "取得に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, &response.DiagnosisDetailResponse{
		Diagnosis: *detail,
	})
}
