package api

import (
	"fmt"
	"path/filepath"

	"development-of-internet-applications/internal/app/handler"
	"development-of-internet-applications/internal/app/repository"
	"development-of-internet-applications/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "development-of-internet-applications/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func StartServer() {
	logrus.Debug("Server started")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория: ", err)
		return
	}

	baseHandler := handler.NewBaseHandler(repo)
	requestContainHandler := handler.NewRequestClassHandler(repo)
	classHandler := handler.NewClassHandler(repo)
	massRequestHandler := handler.NewMassRequestHandler(repo)
	userHandler := handler.NewUserHandler(repo)

	rootDir := filepath.Clean(filepath.Join(filepath.Dir("."), "..", ".."))
	resourcesDir := filepath.Join(rootDir, "resources")
	stylesPath := filepath.Join(resourcesDir, "styles")
	imagesPath := filepath.Join(resourcesDir, "images")

	router := gin.Default()

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware для аутентификации
	authMiddleware := baseHandler.WithAuthCheck(role.User, role.Moderator)

	api := router.Group("/api")
	{
		// Публичные методы (чтение данных) - доступны без авторизации
		api.GET("/classes", classHandler.GetClasses)
		api.GET("/classes/:id", classHandler.GetClassByIDAPI)
		
		// Аутентификация
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		// Защищенные методы - требуют авторизации
		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			// Методы пользователя
			protected.GET("/mass-requests/star-calculation", massRequestHandler.GetStarCalc)        
			protected.GET("/mass-requests", massRequestHandler.GetRequests)           
			protected.POST("/mass-requests", massRequestHandler.CreateRequest)         
			protected.GET("/mass-requests/:id", massRequestHandler.GetRequestByID)       
			protected.PUT("/mass-requests/:id", massRequestHandler.UpdateRequest)     
			protected.PUT("/mass-requests/:id/form", massRequestHandler.FormRequest) 
			protected.DELETE("/mass-requests/:id", massRequestHandler.DeleteRequest)   
			protected.POST("/mass-requests/classes", requestContainHandler.AddClassToRequest) 
			protected.DELETE("/mass-requests/:id/classes/:class_id", requestContainHandler.RemoveClassFromRequest)
			protected.PUT("/mass-requests/:id/classes/:class_id", requestContainHandler.UpdateRequestClassItem)
			protected.POST("/auth/logout", userHandler.Logout)
			protected.GET("/user/profile", userHandler.GetProfile)
			protected.PUT("/user/profile", userHandler.UpdateProfile)

			// Методы модератора (будут проверять роль внутри handler)
			protected.PUT("/mass-requests/:id/complete", massRequestHandler.CompleteRequest)
			protected.POST("/classes", classHandler.CreateClass)
			protected.PUT("/classes/:id", classHandler.UpdateClass)  
			protected.DELETE("/classes/:id", classHandler.DeleteClass)
		}
	}

	router.Static("/resources/styles", stylesPath)
	router.Static("/resources/images", imagesPath)

	fmt.Println("Сервер запущен на http://0.0.0.0:8080")
	router.Run("0.0.0.0:8080")
	logrus.Debug("Server stopped")
}