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

// AddClassToRequest godoc
// @Summary Add class to request
// @Description Add spectral class to draft mass request
// @Tags request-classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]uint64 true "Class data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /mass-requests/classes [post]
func (h *RequestClassHandler) AddClassToRequest(ctx *gin.Context) {
	var req struct {
		ClassID uint64 `json:"class_id" binding:"required"`
	}
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

	if err := h.repo.MassRequest.AddClassToRequest(request.ID, req.ClassID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add class to request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Class added to request successfully"})
}

// RemoveClassFromRequest godoc
// @Summary Remove class from request
// @Description Remove spectral class from mass request
// @Tags request-classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_id path int true "Request ID"
// @Param class_id path int true "Class ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /mass-requests/{request_id}/classes/{class_id} [delete]
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

	// Проверяем что пользователь удаляет класс из своей заявки
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Class removed from request successfully"})
}

// UpdateRequestClassItem godoc
// @Summary Update request class item
// @Description Update class parameters in mass request
// @Tags request-classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_id path int true "Request ID"
// @Param class_id path int true "Class ID"
// @Param updates body map[string]interface{} true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /mass-requests/{id}/classes/{class_id} [put]
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

	// Проверяем что пользователь обновляет свою заявку
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