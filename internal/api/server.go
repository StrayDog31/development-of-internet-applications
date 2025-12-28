package api

import (
	"fmt"
	"path/filepath"
	"time"

	"development-of-internet-applications/internal/app/handler"
	"development-of-internet-applications/internal/app/repository"
	"development-of-internet-applications/internal/app/role"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "development-of-internet-applications/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://172.19.80.1:5173",
			"http://127.0.0.1:5173",
			"http://172.19.80.1:3000",
			"http://192.168.1.100:5173",
			"http://192.168.1.42:1420",
			"https://*.github.io",
			"http://172.19.80.1:1420",
			"http://127.0.0.1:1420",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Api-Key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowWildcard:    true,
		AllowOriginFunc: func(origin string) bool {
			logrus.Debugf("CORS request from origin: %s", origin)
			return true
		},
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authMiddleware := baseHandler.WithAuthCheck(role.User, role.Moderator)

	api := router.Group("/api")
	{
		api.GET("/classes", classHandler.GetClasses)
		api.GET("/classes/:id", classHandler.GetClassByIDAPI)
		
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		api.GET("/mass-requests/star-calculation", massRequestHandler.GetStarCalc)

		protected := api.Group("")
		protected.Use(authMiddleware)
		{
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

			protected.PUT("/mass-requests/:id/complete", massRequestHandler.CompleteRequest)
			protected.POST("/classes", classHandler.CreateClass)
			protected.PUT("/classes/:id", classHandler.UpdateClass)  
			protected.DELETE("/classes/:id", classHandler.DeleteClass)
			protected.PUT("/classes/:id/image", classHandler.UpdateClassImage)
		}
	}

	router.POST("/api/v1/webhook/calculation-result", massRequestHandler.WebhookResult)

	router.Static("/resources/styles", stylesPath)
	router.Static("/resources/images", imagesPath)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	fmt.Println("Сервер запущен на http://0.0.0.0:8080")
	fmt.Println("Swagger доступен по адресу: http://localhost:8080/swagger/index.html")
	router.Run("0.0.0.0:8080")
	logrus.Debug("Server stopped")
}