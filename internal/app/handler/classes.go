package handler

import (
	"development-of-internet-application/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClassesHandler struct {
	ClassRepository       *repository.ClassRepository
	CalcRequestRepository *repository.CalcRequestRepository
}

func NewClassesHandler(classRepo *repository.ClassRepository, calcReqRepo *repository.CalcRequestRepository) *ClassesHandler {
	return &ClassesHandler{
		ClassRepository:       classRepo,
		CalcRequestRepository: calcReqRepo,
	}
}

func (h *ClassesHandler) IndexHandler(ctx *gin.Context) {
	classes, err := h.ClassRepository.GetClasses()
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"cards":     classes, 
	})
}

func (h *ClassesHandler) SearchClasses(ctx *gin.Context) {
	searchQuery := ctx.Query("query")
	var classes []repository.SpectralClass
	var err error

	if searchQuery == "" {
		classes, err = h.ClassRepository.GetClasses()
	} else {
		classes, err = h.ClassRepository.GetClassByName(searchQuery)
	}

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}


	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"cards": classes,
		"searchQuery": searchQuery,
	})
}