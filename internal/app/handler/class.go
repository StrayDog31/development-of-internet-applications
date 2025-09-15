package handler

import (
	"development-of-internet-application/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClassHandler struct {
	ClassRepository *repository.ClassRepository
}

func NewClassHandler(classRepo *repository.ClassRepository) *ClassHandler {
	return &ClassHandler{
		ClassRepository: classRepo,
	}
}

func (handler *ClassHandler) GetClassByID(ctx *gin.Context) {
		idStr := ctx.Param("id")
	ClassID, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	fmt.Sscanf(idStr, "%d", &ClassID)

	class, err := handler.ClassRepository.GetClassByID(ClassID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	
	ctx.HTML(http.StatusOK, "class.html", gin.H{
		"class": class,
	})
}