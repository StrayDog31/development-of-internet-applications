package main


import (
	"development-of-internet-applications/internal/api"
	
	"github.com/sirupsen/logrus"
)

// @title Stellar Mass Calculation API
// @version 1.0
// @description API для расчета масс звезд в заявках
// @host localhost:8080  
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	
	logrus.SetLevel(logrus.DebugLevel)
	api.StartServer()
}