package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type AvatarController struct {
	container *dig.Container
}

func NewAvatarController(container *dig.Container) *AvatarController {
	return &AvatarController{container: container}
}

// @Summary アバター一覧取得（マッチングポイント付き）
// @Tags avatars
// @Description 異性のアバター一覧を取得し、ユーザーとの関係性（マッチングポイント）も含める
// @Security Bearer
// @Success 200 {object} response.AvatarListResponse "アバター一覧取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /avatars [get]
func (c AvatarController) GetAvatarList(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Message: "認証されていない、またはトークンが不正",
		})
	}

	s := service.NewAvatarService(c.container)
	avatars, err := s.GetAvatarList(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, &response.ErrorResponse{
				Message: "ユーザーが見つかりません",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Message: "アバター取得に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, &response.AvatarListResponse{
		Avatars: avatars,
	})
}

// @Summary マッチングポイント更新
// @Tags avatars
// @Description 指定されたアバターとのマッチングポイントを更新する
// @Security Bearer
// @Param avatarId path string true "アバターID（ULID）"
// @Param body body requests.UpdateMatchingPointRequest true "更新するポイント数"
// @Success 200 {object} response.UpdateMatchingPointResponse "マッチングポイント更新成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 404 {object} response.ErrorResponse "アバターが見つからない"
// @Router /avatars/{avatarId}/matching-points [post]
func (c AvatarController) UpdateMatchingPoint(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Message: "認証されていない、またはトークンが不正",
		})
	}
	avatarID := ctx.Param("avatarId")

	var req requests.UpdateMatchingPointRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Message: "リクエストが不正です",
		})
	}

	s := service.NewAvatarService(c.container)
	result, err := s.UpdateMatchingPoint(userID, avatarID, req.Points)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.JSON(http.StatusNotFound, &response.ErrorResponse{
				Message: "アバターが見つかりません",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Message: "マッチングポイントの更新に失敗しました",
		})
	}

	return ctx.JSON(http.StatusOK, result)
}
