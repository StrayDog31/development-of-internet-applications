package repository

import (
	"development-of-internet-applications/internal/app/jwt"
	redisClient "development-of-internet-applications/internal/app/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Class       *ClassRepository
	MassRequest *MassRequestRepository
	User        *UserRepository
	Redis       *redisClient.Client
	Config      *config.Config
}

func NewRepository() (*Repository, error) {
	dsnString := "host=localhost port=5432 user=postgres password=postgres dbname=dia_db sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	cfg := config.LoadConfig()

	// ⚠️ РАСКОММЕНТИРУЙТЕ Redis подключение:
	redis, err := redisClient.New(cfg.Redis)
	if err != nil {
		logrus.Warnf("Failed to connect to Redis: %v", err)
	} else {
		logrus.Info("✅ Redis connected successfully")
	}

	return &Repository{
		Class:       NewClassRepository(db),
		MassRequest: NewMassRequestRepository(db),
		User:        NewUserRepository(db),
		Config:      cfg,
		Redis:       redis, // ⚠️ Теперь Redis будет работать
	}, nil
}