package handler

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestClassHandler struct {
	repo *repository.Repository
}

func NewRequestClassHandler(repository *repository.Repository) *RequestClassHandler {
	return &RequestClassHandler{
		repo: repository,
	}
}

func (h *RequestClassHandler) AddClassToRequest(ctx *gin.Context) {
	var req struct {
		ClassID uint64 `json:"class_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := GetFixedUserID()

	request, err := h.repo.MassRequest.GetUserDraftRequest(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if request == nil {
		request = &ds.MassRequest{
			Status:    1,
			UserID:    userID,
			CreatedAt: time.Now(),
		}
		if err := h.repo.MassRequest.CreateRequest(request); err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
	}

	if err := h.repo.MassRequest.AddClassToRequest(request.ID, req.ClassID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add class to request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Class added to request successfully"})
}

func (h *RequestClassHandler) RemoveClassFromRequest(ctx *gin.Context) {
	requestIDStr := ctx.Param("request_id")
	classIDStr := ctx.Param("class_id")

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	classID, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	if err := h.repo.MassRequest.RemoveClassFromRequest(requestID, classID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove class from request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Class removed from request successfully"})
}

func (h *RequestClassHandler) UpdateRequestClassItem(ctx *gin.Context) {
	requestIDStr := ctx.Param("request_id")
	classIDStr := ctx.Param("class_id")

	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	classID, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	allowedFields := map[string]bool{
		"lyminosity": true,
		"mass":       true,
		"number":     true,
	}

	for field := range updates {
		if !allowedFields[field] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field: " + field})
			return
		}
	}

	if err := h.repo.MassRequest.UpdateRequestClassItem(requestID, classID, updates); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request class item"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request class item updated successfully"})
}
