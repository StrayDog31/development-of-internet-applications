package api

import (
	"fmt"
	"path/filepath"

	"development-of-internet-applications/internal/app/handler"
	"development-of-internet-applications/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	logrus.Debug("Server started")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория: ", err)
		return
	}

	requestContainHandler := handler.NewRequestClassHandler(repo)
	classHandler := handler.NewClassHandler(repo)
	massRequestHandler := handler.NewMassRequestHandler(repo)
	userHandler := handler.NewUserHandler(repo)

	rootDir := filepath.Clean(filepath.Join(filepath.Dir("."), "..", ".."))
	resourcesDir := filepath.Join(rootDir, "resources")
	stylesPath := filepath.Join(resourcesDir, "styles")
	imagesPath := filepath.Join(resourcesDir, "images")

	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/classes", classHandler.GetClasses)
		api.GET("/classes/:id", classHandler.GetClassByIDAPI)
		api.POST("/classes", classHandler.CreateClass)
		api.PUT("/classes/:id", classHandler.UpdateClass)  
		api.DELETE("/classes/:id", classHandler.DeleteClass)

		// Requests API  
		api.GET("/mass-requests/star-calculation", massRequestHandler.GetStarCalc)        
		api.GET("/mass-requests", massRequestHandler.GetRequests)           
		api.POST("/mass-requests", massRequestHandler.CreateRequest)         
		api.GET("/mass-requests/:id", massRequestHandler.GetRequestByID)       
		api.PUT("/mass-requests/:id", massRequestHandler.UpdateRequest)     
		api.PUT("/mass-requests/:id/form", massRequestHandler.FormRequest) 
		api.PUT("/mass-requests/:id/complete", massRequestHandler.CompleteRequest) 
		api.DELETE("/mass-requests/:id", massRequestHandler.DeleteRequest)   

		// Request-Classes API (связи)
		api.POST("/mass-requests/classes", requestContainHandler.AddClassToRequest) 
		api.DELETE("/mass-requests/:id/classes/:class_id", requestContainHandler.RemoveClassFromRequest)
		api.PUT("/mass-requests/:id/classes/:class_id", requestContainHandler.UpdateRequestClassItem)

		// Users API
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)
		api.POST("/auth/logout", userHandler.Logout)
		api.GET("/user/profile", userHandler.GetProfile)
		api.PUT("/user/profile", userHandler.UpdateProfile)
	}

	router.Static("/resources/styles", stylesPath)
	router.Static("/resources/images", imagesPath)

	fmt.Println("Сервер запущен на http://0.0.0.0:8080")
	router.Run("0.0.0.0:8080")
	logrus.Debug("Server stopped")
}