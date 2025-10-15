package handler

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MassRequestHandler struct {
	repo *repository.Repository
}

func NewMassRequestHandler(repository *repository.Repository) *MassRequestHandler {
	return &MassRequestHandler{
		repo: repository,
	}
}

func (h *MassRequestHandler) GetMassRequestByID(ctx *gin.Context) {
	id := ctx.Param("id")
	reqID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	request, err := h.repo.MassRequest.GetMassRequestByID(uint64(reqID), 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	if request == nil || request.ID == 0 {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "request.html", gin.H{
		"request": request,
	})
}

func (h *MassRequestHandler) DeleteRequestHandler(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		logrus.Error("Invalid request ID:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.repo.MassRequest.DeleteRequest(requestID)
	if err != nil {
		logrus.Error("Error deleting request:", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

func (h *MassRequestHandler) GetStarCalc(ctx *gin.Context) {
	userID := GetFixedUserID()

	request, err := h.repo.MassRequest.GetUserDraftRequest(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var requestID uint64
	var itemsCount int64

	if request != nil {
		requestID = request.ID
		itemsCount, err = h.repo.MassRequest.GetRequestItemsCount(requestID)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"request_id":  requestID,
		"items_count": itemsCount,
	})
}

func (h *MassRequestHandler) GetRequests(ctx *gin.Context) {
	filter := &repository.MassRequestFilter{}

	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.ParseUint(statusStr, 10, 8); err == nil {
			filter.Status = uint8(status)
		}
	}

	if startDateStr := ctx.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = startDate
		}
	}

	if endDateStr := ctx.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = endDate
		}
	}

	requests, err := h.repo.MassRequest.GetRequests(filter)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

func (h *MassRequestHandler) GetRequestByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	request, err := h.repo.MassRequest.GetRequestByID(requestID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	ctx.JSON(http.StatusOK, request)
}

func (h *MassRequestHandler) CreateRequest(ctx *gin.Context) {
	userID := GetFixedUserID()

	request := &ds.MassRequest{
		Status:    1,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := h.repo.MassRequest.CreateRequest(request); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Request created successfully",
		"request_id": request.ID,
	})
}

func (h *MassRequestHandler) UpdateRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	request, err := h.repo.MassRequest.GetRequestByID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only draft requests can be modified"})
		return
	}

	if err := h.repo.MassRequest.UpdateRequest(requestID, updates); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request updated successfully"})
}

func (h *MassRequestHandler) FormRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	request, err := h.repo.MassRequest.GetRequestByID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only draft requests can be formed"})
		return
	}

	itemsCount, err := h.repo.MassRequest.GetRequestItemsCount(requestID)
	if err != nil || itemsCount == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request must contain at least one class"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":    3,
		"formed_at": now,
	}

	if err := h.repo.MassRequest.UpdateRequest(requestID, updates); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to form request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request formed successfully"})
}

func (h *MassRequestHandler) CompleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var req struct {
		Action string `json:"action" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	request, err := h.repo.MassRequest.GetRequestByID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only pending requests can be completed/rejected"})
		return
	}

	moderatorID := uint64(2)
	now := time.Now()

	var newStatus uint8
	switch req.Action {
	case "complete":
		newStatus = 4
		if err := h.calculateRequestFields(request); err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate star masses"})
			return
		}
	case "reject":
		newStatus = 5
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	updates := map[string]interface{}{
		"status":       newStatus,
		"moderator_id": moderatorID,
		"closed_at":    now,
	}

	if err := h.repo.MassRequest.UpdateRequest(requestID, updates); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request " + req.Action + "d successfully",
		"status":  newStatus,
	})
}

func (h *MassRequestHandler) DeleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	request, err := h.repo.MassRequest.GetRequestByID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only draft requests can be deleted"})
		return
	}

	if err := h.repo.MassRequest.DeleteRequest(requestID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}

func (h *MassRequestHandler) calculateStarMass(spectralClass string, luminosity uint64) float64 {
	if len(spectralClass) == 0 {
		return 0
	}
	
	class := strings.ToUpper(string(spectralClass[0]))
	
	var exponent float64
	
	switch class {
	case "O", "B", "A":
		exponent = 0.222
	case "F", "G":
		exponent = 0.234
	case "K", "M":
		exponent = 0.264
	default:
		exponent = 0.25
	}
	
	luminosityRatio := float64(luminosity)
	
	mass := math.Pow(luminosityRatio, exponent)
	
	return math.Round(mass*100) / 100
}

func (h *MassRequestHandler) calculateRequestFields(request *ds.MassRequest) error {
	logrus.Infof("Calculating star masses for request ID: %d", request.ID)
	
	for i := range request.MassRequestToClass {
		item := &request.MassRequestToClass[i]
		if item.Luminosity != nil && item.Class.Spectre != "" {
			mass := h.calculateStarMass(item.Class.Spectre, *item.Luminosity)
			massUint := uint64(mass * 100)
			item.Mass = &massUint
			
			logrus.Infof("Calculated mass for class %s (%s): %.2f M☉", 
				item.Class.Name, item.Class.Spectre, mass)
		} else {
			logrus.Warnf("Missing luminosity or spectre for class %s in request %d", 
				item.Class.Name, request.ID)
		}
	}
	
	for _, item := range request.MassRequestToClass {
		if item.Mass != nil {
			updates := map[string]interface{}{
				"mass": item.Mass,
			}
			if err := h.repo.MassRequest.UpdateRequestClassItem(request.ID, item.ClassID, updates); err != nil {
				logrus.Errorf("Failed to update mass for class %d in request %d: %v", 
					item.ClassID, request.ID, err)
				return err
			}
		}
	}
	
	logrus.Infof("Successfully calculated masses for %d classes in request %d", 
		len(request.MassRequestToClass), request.ID)
	return nil
}