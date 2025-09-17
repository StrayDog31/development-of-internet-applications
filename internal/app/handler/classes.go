package handler

import (
	"development-of-internet-application/internal/app/ds"
	"development-of-internet-application/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClassesHandler struct {
	repo *repository.Repository
}

func NewClassesHandler(repository *repository.Repository) *ClassesHandler {
	return &ClassesHandler{
		repo: repository,
	}
}

func (h *ClassesHandler) IndexHandler(ctx *gin.Context) {

	classes, err := h.repo.Class.GetClasses()
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var activeRequestID uint64
	var currentRequestContains int
	activeRequest, err := h.repo.CalcRequest.GetActiveCalcRequestByUser(1)
	if err != nil {
		logrus.Warnf("No active calc request found: %v", err)
		activeRequestID = 0
		currentRequestContains = 0
	} else if activeRequest != nil {
		activeRequestID = activeRequest.ID
		currentRequestContains = len(activeRequest.CalcRequestToClass)
	}
	

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"cards":     classes, 
		"currentRequestID": activeRequestID,
		"currentRequestContains": currentRequestContains,

	})
}

func (h *ClassesHandler) SearchClasses(ctx *gin.Context) {
	searchQuery := ctx.Query("query")
	var classes []ds.Class
	var err error

	if searchQuery == "" {
		classes, err = h.repo.Class.GetClasses()
	} else {
		classes, err = h.repo.Class.GetClassByName(searchQuery)
	}

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var activeRequestID uint64
	var currentRequestContains int
	activeRequest, err := h.repo.CalcRequest.GetActiveCalcRequestByUser(1)
	if err != nil {
		logrus.Warnf("No active calc request found: %v", err)
		activeRequestID = 0
		currentRequestContains = 0
	} else if activeRequest != nil {
		activeRequestID = activeRequest.ID
		currentRequestContains = len(activeRequest.CalcRequestToClass)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"cards": classes,
		"searchQuery": searchQuery,
		"currentRequestID": activeRequestID,
		"currentRequestContains": currentRequestContains,

	})
}

func (h *ClassesHandler) AddClassToRequest(ctx *gin.Context) {
	classIDStr := ctx.PostForm("class-id")
	classID, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.repo.CalcRequest.AddClassToRequest(classID, 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}