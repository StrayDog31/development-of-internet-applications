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

type AddClassRequest struct {
	ClassID    uint64  `json:"class_id" binding:"required"`
	Luminosity *uint64 `json:"luminosity,omitempty"`
	Mass       *uint64 `json:"mass,omitempty"`
}

func (h *RequestClassHandler) AddClassToRequest(ctx *gin.Context) {
	var req AddClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	if err := h.repo.MassRequest.AddClassToRequest(request.ID, req.ClassID, req.Luminosity, req.Mass); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add class to request"})
		return
	}

	updatedRequest, err := h.repo.MassRequest.GetMassRequestByID(request.ID, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated request"})
		return
	}

	ctx.JSON(http.StatusOK, updatedRequest)
}

func (h *RequestClassHandler) RemoveClassFromRequest(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
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

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	request, err := h.repo.MassRequest.GetMassRequestByID(requestID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only draft requests can be modified"})
		return
	}

	if err := h.repo.MassRequest.RemoveClassFromRequest(requestID, classID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove class from request"})
		return
	}

	updatedRequest, err := h.repo.MassRequest.GetMassRequestByID(requestID, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated request"})
		return
	}

	ctx.JSON(http.StatusOK, updatedRequest)
}

func (h *RequestClassHandler) UpdateRequestClassItem(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
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

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	request, err := h.repo.MassRequest.GetMassRequestByID(requestID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only draft requests can be modified"})
		return
	}

	allowedFields := map[string]bool{
		"luminosity": true,
		"mass":       true,
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

	updatedRequest, err := h.repo.MassRequest.GetMassRequestByID(requestID, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated request"})
		return
	}

	ctx.JSON(http.StatusOK, updatedRequest)
}