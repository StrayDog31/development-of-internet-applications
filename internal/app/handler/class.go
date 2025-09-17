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
	repo *repository.Repository
}

func NewClassHandler(repository *repository.Repository) *ClassHandler {
	return &ClassHandler{
		repo: repository,
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

	class, err := handler.repo.Class.GetClassByID(uint64(ClassID))
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}
	
	ctx.HTML(http.StatusOK, "class.html", gin.H{
		"class": class,
	})
}