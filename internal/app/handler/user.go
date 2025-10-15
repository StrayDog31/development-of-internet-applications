package handler

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"development-of-internet-applications/internal/app/role"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  uint64 `json:"user_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	UserID      uint64 `json:"user_id"`
	Role        string `json:"role"`
}

type UpdateProfileRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// @Summary      Register a new user
// @Description  Create a new user account with username and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "User registration data"
// @Success      201  {object}  RegisterResponse
// @Failure      400  {object}  map[string]string
// @Router       /users/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	if req.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("Failed to hash password:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user := &ds.User{
		Username: req.Username,
		PassHash: string(hashedPassword),
		IsMod:    false,
	}

	if err := h.repo.User.CreateUser(user); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, RegisterResponse{
		Message: "User registered successfully",
		UserID:  user.ID,
	})
}

// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "User credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /users/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logrus.Error("Bind JSON error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	logrus.Debugf("Login attempt for user: %s", req.Username)

	user, err := h.repo.User.GetUserByUsername(req.Username)
	if err != nil {
		logrus.Errorf("User not found: %s, error: %v", req.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	logrus.Debugf("User found: %s, ID: %d, IsMod: %t", user.Username, user.ID, user.IsMod)

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(req.Password))
	if err != nil {
		logrus.Errorf("Password mismatch for user: %s, error: %v", req.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	logrus.Debug("Password verified successfully")

	// Проверка конфигурации JWT
	if h.repo == nil {
		logrus.Error("Repository is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error: repository not initialized"})
		return
	}
	if h.repo.Config == nil {
		logrus.Error("Config is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error: config not loaded"})
		return
	}

	// JWTConfig - это структура, проверяем на пустоту вместо nil
	cfg := h.repo.Config.JWT
	if cfg.Token == "" {
		logrus.Error("JWT token is empty in config")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error: JWT token not configured"})
		return
	}

	logrus.Debugf("JWT config loaded: expires_in=%v", cfg.ExpiresIn)

	expiresAt := time.Now().Add(cfg.ExpiresIn)

	// Конвертируем IsMod в роль
	var userRole role.Role
	if user.IsMod {
		userRole = role.Moderator
	} else {
		userRole = role.User
	}

	logrus.Debugf("User role: %s", userRole.String())

	token := jwt.NewWithClaims(cfg.SigningMethod, &ds.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "stellar-backend",
		},
		UserID: user.ID,
		Role:   userRole,
	})

	if token == nil {
		logrus.Error("Failed to create JWT token - token is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	strToken, err := token.SignedString([]byte(cfg.Token))
	if err != nil {
		logrus.Errorf("Failed to sign JWT token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	logrus.Debugf("JWT token created successfully for user: %s", user.Username)

	response := LoginResponse{
		ExpiresIn:   int64(cfg.ExpiresIn.Seconds()),
		AccessToken: strToken,
		TokenType:   "Bearer",
		UserID:      user.ID,
		Role:        userRole.String(),
	}

	logrus.Debugf("Login successful for user: %s, user_id: %d, role: %s", 
		user.Username, user.ID, userRole.String())

	ctx.JSON(http.StatusOK, response)
}

// @Summary      User logout
// @Description  Invalidate JWT token by adding it to blacklist
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization header"})
		return
	}

	jwtStr = jwtStr[len(jwtPrefix):]

	token, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.repo.Config.JWT.Token), nil
	})
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*ds.JWTClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
		return
	}

	if h.repo.Redis != nil {
		err = h.repo.Redis.WriteJWTToBlacklist(ctx.Request.Context(), jwtStr, ttl)
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// @Summary      Get user profile
// @Description  Get current authenticated user's profile information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  ds.User
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.repo.User.GetUserByID(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Не возвращаем хеш пароля
	user.PassHash = ""
	ctx.JSON(http.StatusOK, user)
}

// @Summary      Update user profile
// @Description  Update current authenticated user's profile (username and/or password)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateProfileRequest true "Profile update data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /users/profile [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Хешируем новый пароль если он предоставлен
	var hashedPassword *string
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Error("Failed to hash password:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
		hashedStr := string(hashed)
		hashedPassword = &hashedStr
	}

	if err := h.repo.User.UpdateUser(userID, req.Username, hashedPassword); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}