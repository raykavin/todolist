package handler

import (
	"errors"
	netHttp "net/http"

	"todolist/internal/adapter/delivery/http"
	entPerson "todolist/internal/domain/person/entity"
	"todolist/internal/domain/person/valueobject"
	entUser "todolist/internal/domain/user/entity"
	"todolist/internal/dto"
	ucPerson "todolist/internal/usecase/person"
	ucUser "todolist/internal/usecase/user"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	createUserUseCase     ucUser.CreateUserUseCase
	createPersonUseCase   ucPerson.CreatePersonUseCase
	loginUseCase          ucUser.LoginUseCase
	changePasswordUseCase ucUser.ChangePasswordUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	createUserUseCase ucUser.CreateUserUseCase,
	createPersonUseCase ucPerson.CreatePersonUseCase,
	loginUseCase ucUser.LoginUseCase,
	changePasswordUseCase ucUser.ChangePasswordUseCase,
) *AuthHandler {
	return &AuthHandler{
		createUserUseCase:     createUserUseCase,
		createPersonUseCase:   createPersonUseCase,
		loginUseCase:          loginUseCase,
		changePasswordUseCase: changePasswordUseCase,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} dto.Response{data=dto.UserResponse}
// @Failure 400 {object} dto.Response
// @Failure 409 {object} dto.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(ctx http.RequestContext) {
	var input dto.CreateUserRequest

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))
		ctx.Abort()
		return
	}

	person, err := h.createPersonUseCase.Execute(ctx.Context(), input.Person)
	if err != nil {
		if errors.Is(err, valueobject.ErrInvalidEmail) {
			ctx.JSON(netHttp.StatusConflict,
				dto.ErrorResponse("EMAIL_EXISTS", "Email already exists", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("CREATE_FAILED", "Failed to create person", nil))
		}

		ctx.Abort()
		return
	}

	// Create user
	user, err := h.createUserUseCase.Execute(ctx.Context(), person.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, entPerson.ErrPersonNotFound):
			ctx.JSON(netHttp.StatusBadRequest,
				dto.ErrorResponse("PERSON_NOT_FOUND", "Person not found", nil))
		case errors.Is(err, entUser.ErrUsernameExists):
			ctx.JSON(netHttp.StatusConflict,
				dto.ErrorResponse("USERNAME_EXISTS", "Username already exists", nil))
		case errors.Is(err, entUser.ErrUserAlreadyExists):
			ctx.JSON(netHttp.StatusConflict,
				dto.ErrorResponse("USER_EXISTS", "User already exists for this person", nil))
		default:
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("REGISTRATION_FAILED", "Failed to register user", nil))
		}
		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusCreated,
		dto.SuccessResponse(user, "User registered successfully"))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.AuthRequest true "Login credentials"
// @Success 200 {object} dto.Response{data=dto.AuthResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx http.RequestContext) {
	var input dto.AuthRequest

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))
		ctx.Abort()
		return
	}

	// Authenticate user
	authResponse, err := h.loginUseCase.Execute(ctx.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ucUser.ErrInvalidCredentials):
			ctx.JSON(netHttp.StatusUnauthorized,
				dto.ErrorResponse("INVALID_CREDENTIALS", "Invalid username or password", nil))
		case errors.Is(err, ucUser.ErrUserNotActive):
			ctx.JSON(netHttp.StatusUnauthorized,
				dto.ErrorResponse("USER_INACTIVE", "User account is not active", nil))
		default:
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("LOGIN_FAILED", "Failed to login", nil))
		}
		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(authResponse, "Login successful"))
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Param passwords body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Security BearerAuth
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(ctx http.RequestContext) {
	userID, err := getAuthenticatedUserID(ctx)
	if err != nil || userID == 0 {
		ctx.JSON(netHttp.StatusUnauthorized,
			dto.ErrorResponse("UNAUTHENTICATED", "User is not authenticated", nil))
		ctx.Abort()
		return
	}

	var input dto.ChangePasswordRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(netHttp.StatusBadRequest,
			dto.ErrorResponse("INVALID_REQUEST", "Invalid request body", parseError(err)))
		ctx.Abort()
		return
	}

	// Change password
	err = h.changePasswordUseCase.Execute(ctx.Context(), userID, input)
	if err != nil {
		if err.Error() == "old password is incorrect" {
			ctx.JSON(netHttp.StatusBadRequest,
				dto.ErrorResponse("INCORRECT_PASSWORD", "Old password is incorrect", nil))
		} else {
			ctx.JSON(netHttp.StatusInternalServerError,
				dto.ErrorResponse("CHANGE_PASSWORD_FAILED", "Failed to change password", nil))
		}
		ctx.Abort()
		return
	}

	ctx.JSON(netHttp.StatusOK, dto.SuccessResponse(nil, "Password changed successfully"))
}

// Register registers auth routes
// Example usage with Gin:
// func (h *AuthHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
// 	auth := router.Group("/auth")
// 	{
// 		// Public routes
// 		auth.POST("/register", h.Register)
// 		auth.POST("/login", h.Login)
//
// 		// Protected routes
// 		auth.Use(authMiddleware)
// 		auth.POST("/change-password", h.ChangePassword)
// 		auth.GET("/me", h.GetCurrentUser)
// 	}
// }
