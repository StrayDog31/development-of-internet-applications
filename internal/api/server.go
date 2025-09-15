package api

import (
	"fmt"
	"path/filepath"

	"development-of-internet-application/internal/app/handler"
	"development-of-internet-application/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	logrus.Debug("Server started")

	classRepo, err := repository.NewClassRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория приборов")
	}
	requestRepo, err := repository.NewCalcRequestRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория запросов")
	}

	calcRequestHandler := handler.NewCalcRequestHandler(classRepo, requestRepo)
	classHandler := handler.NewClassHandler(classRepo)
	classesHandler := handler.NewClassesHandler(classRepo, requestRepo)

	rootDir := filepath.Clean(filepath.Join(filepath.Dir("."), "..", ".."))
	
	resourcesDir := filepath.Join(rootDir, "resources")
	stylesPath := filepath.Join(resourcesDir, "styles")
	imagesPath := filepath.Join(resourcesDir, "images")
	templatesPath := filepath.Join(rootDir, "templates")

	router := gin.Default()

	router.LoadHTMLGlob(filepath.Join(templatesPath, "*.html"))

	router.GET("/", classesHandler.IndexHandler)
	router.GET("/class/:id", classHandler.GetClassByID)
	router.GET("/request/:id", calcRequestHandler.GetCalcRequestByID)

	router.GET("/search", classesHandler.SearchClasses)


	router.Static("/resources/styles", stylesPath)
	router.Static("/resources/images", imagesPath)

	fmt.Println("Сервер запущен на http://0.0.0.0:8080")
	router.Run("0.0.0.0:8080")

	logrus.Debug("Server stopped")
}