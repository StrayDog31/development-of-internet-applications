package handler

import (
	"bytes"
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"development-of-internet-applications/internal/app/role"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NestJSResponse struct {
	RequestID  uint64        `json:"request_id"`
	Status     string        `json:"status"`
	Calculated bool          `json:"calculated"`
	Results    []ClassResult `json:"results,omitempty"`
	Error      string        `json:"error,omitempty"`
}

type ClassResult struct {
	ClassID uint64  `json:"class_id"`
	Mass    float64 `json:"mass"`
	Message string  `json:"message,omitempty"`
}

type MassRequestHandler struct {
	repo             *repository.Repository
	webhookBaseURL   string
	nestJSServiceURL string
}

func NewMassRequestHandler(repository *repository.Repository) *MassRequestHandler {
	return &MassRequestHandler{
		repo:             repository,
		webhookBaseURL:   "http://172.19.80.1:8080",
		nestJSServiceURL: "http://172.19.80.1:8000",
	}
}

func (h *MassRequestHandler) GetStarCalc(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"request_id":  -1,
		"items_count": 0,
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

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userRole, exists := GetUserRoleFromContext(ctx)
	if exists && userRole == role.Moderator {
		requests, err := h.repo.MassRequest.GetRequests(filter)
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, requests)
	} else {
		requests, err := h.repo.MassRequest.GetMassRequests(userID, filter)
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		ctx.JSON(http.StatusOK, requests)
	}
}

func (h *MassRequestHandler) GetRequestByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userRole, exists := GetUserRoleFromContext(ctx)

	var request *ds.MassRequest
	if exists && userRole == role.Moderator {
		request, err = h.repo.MassRequest.GetRequestByID(requestID)
	} else {
		request, err = h.repo.MassRequest.GetMassRequestByID(requestID, userID)
	}

	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	ctx.JSON(http.StatusOK, request)
}

func (h *MassRequestHandler) CreateRequest(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	userRole, exists := GetUserRoleFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
		return
	}

	if userRole != role.Moderator {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":     "Moderator access required",
			"your_role": userRole,
		})
		return
	}

	request, err := h.repo.MassRequest.GetRequestByID(requestID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	if request.Status != 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":          "Only pending requests can be completed/rejected",
			"current_status": request.Status,
		})
		return
	}

	moderatorID, _ := GetUserIDFromContext(ctx)
	now := time.Now()

	if req.Action == "complete" {
		success := h.sendToNestJS(request)

		if !success {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send request to calculation service",
			})
			return
		}

		updates := map[string]interface{}{
			"status":       6,
			"moderator_id": moderatorID,
			"closed_at":    now,
		}

		if err := h.repo.MassRequest.UpdateRequest(requestID, updates); err != nil {
			logrus.Errorf("Failed to update request status: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update request status",
			})
			return
		}

		ctx.JSON(http.StatusAccepted, gin.H{
			"message":    "Request accepted for processing",
			"request_id": requestID,
			"status":     6,
		})
		return
	}

	if req.Action == "reject" {
		updates := map[string]interface{}{
			"status":       5,
			"moderator_id": moderatorID,
			"closed_at":    now,
		}
		if err := h.repo.MassRequest.UpdateRequest(requestID, updates); err != nil {
			logrus.Errorf("Failed to reject request: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to reject request",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Request rejected successfully",
			"status":  5,
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
}

func (h *MassRequestHandler) DeleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
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

func (h *MassRequestHandler) sendToNestJS(request *ds.MassRequest) bool {
	type ClassData struct {
		ClassID    uint64  `json:"class_id"`
		Luminosity *uint64 `json:"luminosity,omitempty"`
	}

	type Payload struct {
		RequestID   uint64      `json:"request_id"`
		UserID      uint64      `json:"user_id"`
		CallbackURL string      `json:"callback_url"`
		Classes     []ClassData `json:"classes"`
	}

	var classes []ClassData
	for _, item := range request.MassRequestToClass {
		classes = append(classes, ClassData{
			ClassID:    item.ClassID,
			Luminosity: item.Luminosity,
		})
	}

	if len(classes) == 0 {
		logrus.Errorf("No classes found for request %d", request.ID)
		return false
	}

	callbackURL := fmt.Sprintf("%s/api/v1/webhook/calculation-result", h.webhookBaseURL)

	payload := Payload{
		RequestID:   request.ID,
		UserID:      request.UserID,
		CallbackURL: callbackURL,
		Classes:     classes,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Failed to marshal payload for request %d: %v", request.ID, err)
		return false
	}

	url := fmt.Sprintf("%s/calculator/star-mass", h.nestJSServiceURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		logrus.Errorf("Failed to send request %d to NestJS service (%s): %v", request.ID, url, err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response from NestJS for request %d: %v", request.ID, err)
		return false
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		logrus.Errorf("NestJS returned error for request %d: HTTP %d - %s", request.ID, resp.StatusCode, string(body))
		return false
	}

	var nestResponse struct {
		Message   string `json:"message"`
		RequestID uint64 `json:"request_id"`
		Status    string `json:"status"`
	}

	if err := json.Unmarshal(body, &nestResponse); err != nil {
		logrus.Errorf("Failed to parse NestJS response for request %d: %v", request.ID, err)
	}

	logrus.Infof("Successfully sent request %d to NestJS. Response: %s", request.ID, nestResponse.Message)

	return true
}

func (h *MassRequestHandler) WebhookResult(ctx *gin.Context) {
    apiKey := ctx.GetHeader("X-Api-Key")
    if apiKey != "SECRET_KEY_123" {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
        return
    }

    var result NestJSResponse
    if err := ctx.ShouldBindJSON(&result); err != nil {
        logrus.Errorf("Invalid webhook payload: %v", err)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
        return
    }

    logrus.Infof("Received webhook for request %d with status: %s", result.RequestID, result.Status)

    ctx.JSON(http.StatusAccepted, gin.H{
        "message": "Webhook accepted, processing started",
        "request_id": result.RequestID,
    })

    go h.processWebhookAsync(result)
}

func (h *MassRequestHandler) processWebhookAsync(result NestJSResponse) {
    _, err := h.repo.MassRequest.GetRequestByID(result.RequestID)
    if err != nil {
        logrus.Errorf("Request not found: %d", result.RequestID)
        return
    }

    var newStatus uint8
    var message string

    if result.Calculated && len(result.Results) > 0 {
        successCount := 0
        for _, classResult := range result.Results {
            massUint := uint64(classResult.Mass * 100)
            updates := map[string]interface{}{
                "mass": &massUint,
            }

            if err := h.repo.MassRequest.UpdateRequestClassItem(
                result.RequestID,
                classResult.ClassID,
                updates,
            ); err != nil {
                logrus.Errorf("Failed to update mass for class %d in request %d: %v",
                    classResult.ClassID, result.RequestID, err)
            } else {
                successCount++
            }
        }

        newStatus = 4
        message = fmt.Sprintf("Calculation completed successfully. Updated %d/%d classes",
            successCount, len(result.Results))

        logrus.Infof("Updated masses for %d classes in request %d", successCount, result.RequestID)

    } else {
        newStatus = 7
        if result.Error != "" {
            message = fmt.Sprintf("Calculation failed: %s", result.Error)
        } else {
            message = "Calculation failed: no results returned"
        }
        logrus.Errorf("Calculation failed for request %d: %s", result.RequestID, message)
    }

    updates := map[string]interface{}{
        "status": newStatus,
    }

    if err := h.repo.MassRequest.UpdateRequest(result.RequestID, updates); err != nil {
        logrus.Errorf("Failed to update request status: %v", err)
        return
    }

    logrus.Infof("Request %d updated to status %d: %s", result.RequestID, newStatus, message)
}