package handler

import (
	"development-of-internet-application/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CalcRequestHandler struct {
	repo *repository.Repository
}

func NewCalcRequestHandler(repository *repository.Repository) *CalcRequestHandler {
	return &CalcRequestHandler{
		repo: repository,
	}
}


func (handler *CalcRequestHandler) GetCalcRequestByID(ctx *gin.Context) {
	id := ctx.Param("id")
	reqID, err := strconv.Atoi(id)
	if err != nil {
		logrus.Error(err)
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
	}
	fmt.Printf("Request ID: %s\n", id)

	request, err := handler.repo.CalcRequest.GetCalcRequestByID(uint64(reqID), 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
	}

	if request == nil || request.ID == 0 {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "request.html", gin.H{
		"request": request,
	})
}

func (h *CalcRequestHandler) DeleteRequestHandler(ctx *gin.Context) {
	requestIDStr := ctx.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		logrus.Error("Invalid request ID:", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	userID := uint64(1)

	err = h.repo.CalcRequest.DeleteRequest(requestID, userID)
	if err != nil {
		logrus.Error("Error deleting request:", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}