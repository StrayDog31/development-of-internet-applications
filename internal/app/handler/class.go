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