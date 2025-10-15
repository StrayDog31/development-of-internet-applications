package main

import (
	"development-of-internet-applications/internal/api"
	
	"github.com/sirupsen/logrus"
)

func main() {
	
	logrus.SetLevel(logrus.ErrorLevel)
	api.StartServer()
}