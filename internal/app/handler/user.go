package handler

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var fixedUserID uint64 = 2

func GetFixedUserID() uint64 {
	return fixedUserID
}

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var user ds.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.repo.User.CreateUser(&user); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
	})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userID := GetFixedUserID()

	user, err := h.repo.User.GetUserByID(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := GetFixedUserID()

	var username *string
	var password *string
	
	if val, ok := updates["username"]; ok {
		if str, ok := val.(string); ok {
			username = &str
		}
	}
	if val, ok := updates["password"]; ok {
		if str, ok := val.(string); ok {
			password = &str
		}
	}

	if err := h.repo.User.UpdateUser(userID, username, password); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := h.repo.User.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user_id": user.ID,
		"is_mod":  user.IsMod,
	})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}