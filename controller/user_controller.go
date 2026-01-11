package controller

import (
	"net/http"

	"github.com/hackathon-20260110/api/middleware"
	"github.com/hackathon-20260110/api/requests"
	"github.com/hackathon-20260110/api/response"
	"github.com/hackathon-20260110/api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type UserController struct {
	container *dig.Container
}

func NewUserController(container *dig.Container) *UserController {
	return &UserController{container: container}
}

// @Summary 自分の情報取得（マイページ用）
// @Tags users
// @Description 自分の全情報を取得する（マイページ用）。ユーザーが登録されていない場合はis_registered=falseを返す。
// @Security Bearer
// @Success 200 {object} response.GetMeResponse "自分の情報取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /users/me [get]
func (c *UserController) GetMe(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	s := service.NewUserService(c.container)
	u, err := s.GetUserByIDOrNil(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "ユーザー情報の取得に失敗しました",
		})
	}

	if u == nil {
		return ctx.JSON(http.StatusOK, &response.GetMeResponse{
			IsRegistered: false,
		})
	}

	return ctx.JSON(http.StatusOK, &response.GetMeResponse{
		IsRegistered: true,
		User:         u,
	})
}

// @Summary ユーザー作成
// @Tags users
// @Description ユーザーを作成する
// @Security Bearer
// @Param request body requests.CreateUserRequest true "ユーザー作成リクエスト"
// @Success 200 {object} response.User "ユーザー作成成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /users [post]
func (c *UserController) CreateUser(ctx echo.Context) error {
	userID, ok := ctx.Get("userID").(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, &response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証されていない、またはトークンが不正",
		})
	}

	args := new(requests.CreateUserRequest)
	if err := ctx.Bind(args); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
			Error:   "invalid_request",
			Message: "リクエストが不正です",
		})
	}

	s := service.NewUserService(c.container)
	u, err := s.UpsertUser(userID, *args)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response.ErrorResponse{
			Error:   "internal_server_error",
			Message: "ユーザー情報の更新に失敗しました",
		})
	}
	return ctx.JSON(http.StatusOK, u)
}

// @Summary 自分の情報更新
// @Tags users
// @Description 自分のプロフィール情報を更新する
// @Security Bearer
// @Param request body requests.UpdateUserRequest true "ユーザー更新リクエスト"
// @Success 200 {object} response.User "情報更新成功"
// @Failure 400 {object} response.ErrorResponse "リクエストが不正"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /users/me [put]
func (c *UserController) UpdateMe(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)

	// TODO: リクエストボディをパース
	// var req requests.UpdateUserRequest
	// if err := ctx.Bind(&req); err != nil {
	//     return ctx.JSON(http.StatusBadRequest, &response.ErrorResponse{
	//         Error:   "invalid_request",
	//         Message: "リクエストが不正です",
	//     })
	// }

	// TODO: バリデーション
	// TODO: ユーザー情報を更新
	// TODO: 更新後のユーザー情報を取得

	_ = userID // TODO: ユーザー情報更新で使用

	return ctx.JSON(http.StatusOK, &response.User{
		ID:                  "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		DisplayName:         "山田太郎",
		Age:                 25,
		Gender:              "male",
		ProfileImageURL:     "https://example.com/images/profile.jpg",
		Bio:                 "よろしくお願いします！",
		OnboardingCompleted: true,
		CreatedAt:           "2024-01-01T00:00:00Z",
		UpdatedAt:           "2024-01-01T00:00:00Z",
	})
}

// @Summary 特定ユーザーの公開情報取得
// @Tags users
// @Description 特定ユーザーの公開情報を取得する（段階的に情報が公開される）
// @Security Bearer
// @Param userId path string true "ユーザーID（ULID）"
// @Success 200 {object} response.GetUserDetailResponse "ユーザー情報取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 404 {object} response.ErrorResponse "ユーザーが見つからない"
// @Router /users/{userId} [get]
func (c *UserController) GetUser(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	currentUserID := middleware.GetFirebaseUID(ctx)
	userID := ctx.Param("userId")

	// TODO: ユーザーが存在するか確認
	// TODO: チャット開始済みかどうかを確認（段階的公開のため）
	// TODO: 公開情報を取得（チャット開始状況に応じて情報量を調整）

	_ = currentUserID // TODO: 権限チェックで使用

	return ctx.JSON(http.StatusOK, &response.GetUserDetailResponse{
		User: response.UserDetail{
			ID:                 userID,
			DisplayName:        "佐藤花子",
			Age:                24,
			Gender:             "female",
			ProfileImageURL:    "https://example.com/images/profile2.jpg",
			Bio:                "よろしくお願いします！",
			Interests:          []string{"読書", "映画"},
			Location:           "東京都",
			Occupation:         "エンジニア",
			HasAvatarAI:        true,
			AvatarAIAccessible: false,
			CreatedAt:          "2024-01-01T00:00:00Z",
			UpdatedAt:          "2024-01-01T00:00:00Z",
		},
	})
}

// @Summary 特定ユーザーの分身AI取得
// @Tags users
// @Description 特定ユーザーの分身AI情報を取得する（チャット開始後のみアクセス可能）
// @Security Bearer
// @Param userId path string true "ユーザーID（ULID）"
// @Success 200 {object} response.GetUserAvatarAIResponse "分身AI情報取得成功"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Failure 403 {object} response.ErrorResponse "チャットが開始されていないためアクセス不可"
// @Failure 404 {object} response.ErrorResponse "ユーザーまたは分身AIが見つからない"
// @Router /users/{userId}/avatar-ai [get]
func (c *UserController) GetUserAvatarAI(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	currentUserID := middleware.GetFirebaseUID(ctx)
	userID := ctx.Param("userId")

	// TODO: ユーザーが存在するか確認
	// TODO: チャットが開始されているか確認
	// TODO: 分身AI情報を取得

	_ = currentUserID // TODO: 権限チェックで使用

	return ctx.JSON(http.StatusOK, &response.GetUserAvatarAIResponse{
		AvatarAI: &response.AvatarAI{
			ID:     "01ARZ3NDEKTSV4RRFFQ69G5FDV",
			UserID: userID,
			Name:   "佐藤花子の分身AI",
			Personality: map[string]string{
				"tone":  "gentle",
				"style": "polite",
			},
			Bio:       "佐藤花子の性格や好みを反映した分身AIです",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
	})
}

// @Summary 自分の分身AI情報取得
// @Tags users
// @Description 自分の分身AI情報を取得する
// @Security Bearer
// @Success 200 {object} response.GetAvatarAIResponse "分身AI情報取得成功"
// @Failure 404 {object} response.ErrorResponse "分身AIが存在しない（オンボーディング未完了）"
// @Failure 401 {object} response.ErrorResponse "認証されていない、またはトークンが不正"
// @Router /users/me/avatar-ai [get]
func (c *UserController) GetMyAvatarAI(ctx echo.Context) error {
	// ミドルウェアで検証済みのFirebase UIDを取得
	userID := middleware.GetFirebaseUID(ctx)

	// TODO: ユーザーIDから分身AI情報を取得

	return ctx.JSON(http.StatusOK, &response.GetAvatarAIResponse{
		AvatarAI: &response.AvatarAI{
			ID:     "01ARZ3NDEKTSV4RRFFQ69G5FDV",
			UserID: userID,
			Name:   "私の分身AI",
			Personality: map[string]string{
				"tone":  "friendly",
				"style": "casual",
			},
			Bio:       "あなたの性格や好みを反映した分身AIです",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
	})
}
