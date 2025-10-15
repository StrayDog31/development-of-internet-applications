package main

import (
	"development-of-internet-applications/internal/app/ds"
	"development-of-internet-applications/internal/app/dsn"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Error getting current directory:", err)
	}
	rootDir := filepath.Dir(filepath.Dir(projectRoot))
	envPath := filepath.Join(rootDir, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.Class{},
		&ds.MassRequest{},
		&ds.MassRequestToClass{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}

		err = db.Migrator().DropTable(&ds.MassRequest{})
	if err != nil {
		log.Fatal("Error dropping table:", err)
	}

	err = db.AutoMigrate(&ds.MassRequest{})
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}