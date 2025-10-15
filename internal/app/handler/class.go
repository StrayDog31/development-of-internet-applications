package handler

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClassHandler struct {
	repo *repository.Repository
}

func NewClassHandler(repository *repository.Repository) *ClassHandler {
	return &ClassHandler{
		repo: repository,
	}
}

// GetClassByID godoc
// @Summary Get class page (HTML)
// @Description Get HTML page with class information
// @Tags classes
// @Produce html
// @Param id path int true "Class ID"
// @Success 200 {string} string "HTML page"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Router /classes/{id}/page [get]
func (h *ClassHandler) GetClassByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	classID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logrus.Error("Invalid class ID:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	class, err := h.repo.Class.GetClassByID(classID)
	if err != nil {
		logrus.Error("Error getting class:", err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "class.html", gin.H{
		"class": class,
	})
}

// GetClasses godoc
// @Summary Get all spectral classes
// @Description Get list of spectral classes with optional filtering
// @Tags classes
// @Accept json
// @Produce json
// @Param name query string false "Filter by class name"
// @Param color query string false "Filter by color"
// @Param spectre query string false "Filter by spectral class"
// @Param min_temp query int false "Filter by minimum temperature"
// @Param max_temp query int false "Filter by maximum temperature"
// @Success 200 {array} ds.Class
// @Failure 500 {object} map[string]string
// @Router /classes [get]
func (h *ClassHandler) GetClasses(ctx *gin.Context) {
	filter := &repository.ClassFilter{
		Name:    ctx.Query("name"),
		Color:   ctx.Query("color"),
		Spectre: ctx.Query("spectre"),
	}

	if minTemp := ctx.Query("min_temp"); minTemp != "" {
		if temp, err := strconv.Atoi(minTemp); err == nil {
			filter.MinTemp = temp
		}
	}
	if maxTemp := ctx.Query("max_temp"); maxTemp != "" {
		if temp, err := strconv.Atoi(maxTemp); err == nil {
			filter.MaxTemp = temp
		}
	}

	classes, err := h.repo.Class.GetClasses(filter)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, classes)
}

// CreateClass godoc
// @Summary Create new spectral class
// @Description Create a new spectral class (Moderator only)
// @Tags classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param class body ds.Class true "Class data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /classes [post]
func (h *ClassHandler) CreateClass(ctx *gin.Context) {
	var class ds.Class
	if err := ctx.ShouldBindJSON(&class); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	class.IsDeleted = false

	if err := h.repo.Class.CreateClass(&class); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create class"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Class created successfully",
		"id":      class.ID,
	})
}

// UpdateClass godoc
// @Summary Update spectral class
// @Description Update spectral class information (Moderator only)
// @Tags classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class ID"
// @Param updates body map[string]interface{} true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /classes/{id} [put]
func (h *ClassHandler) UpdateClass(ctx *gin.Context) {
	idStr := ctx.Param("id")
	classID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.repo.Class.UpdateClass(classID, updates); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update class"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Class updated successfully"})
}

// DeleteClass godoc
// @Summary Delete spectral class
// @Description Delete a spectral class (Moderator only)
// @Tags classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /classes/{id} [delete]
func (h *ClassHandler) DeleteClass(ctx *gin.Context) {
	idStr := ctx.Param("id")
	classID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	if err := h.repo.Class.DeleteClass(classID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete class"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}

// GetClassByIDAPI godoc
// @Summary Get spectral class by ID
// @Description Get detailed information about a specific spectral class
// @Tags classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID"
// @Success 200 {object} ds.Class
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /classes/{id} [get]
func (h *ClassHandler) GetClassByIDAPI(ctx *gin.Context) {
	idStr := ctx.Param("id")
	classID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	class, err := h.repo.Class.GetClassByID(classID)
	if err != nil {
		logrus.Error("Error getting class:", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	ctx.JSON(http.StatusOK, class)
}

// UpdateClassImage godoc
// @Summary Update class image URL
// @Description Update spectral class image URL from MinIO
// @Tags classes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class ID"
// @Param request body map[string]string true "Image URL data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /classes/{id}/image [put]
func (h *ClassHandler) UpdateClassImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	classID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid class ID"})
		return
	}

	var req struct {
		ImageURL string `json:"image_url" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	updates := map[string]interface{}{
		"image": req.ImageURL,
	}
	
	if err := h.repo.Class.UpdateClass(classID, updates); err != nil {
		logrus.Error("Failed to update class with image URL:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update class"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Image URL updated successfully",
		"image_url": req.ImageURL,
	})
}
