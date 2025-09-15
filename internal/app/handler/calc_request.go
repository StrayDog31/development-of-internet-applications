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
	ClassRepository *repository.ClassRepository
	CalcRequestRepository *repository.CalcRequestRepository
}

func NewCalcRequestHandler(classRepo *repository.ClassRepository, calcReqRepo *repository.CalcRequestRepository) *CalcRequestHandler {
	return &CalcRequestHandler{
		ClassRepository:       classRepo,
		CalcRequestRepository: calcReqRepo,
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

	request, err := handler.CalcRequestRepository.GetCalcRequestViewByID(reqID, handler.ClassRepository)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
	}
	fmt.Printf("Request stars count: %d\n", len(request.Stars))

	ctx.HTML(http.StatusOK, "request.html", gin.H{
		"request": request,
	})
}