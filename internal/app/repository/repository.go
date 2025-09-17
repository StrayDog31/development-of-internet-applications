package repository

import (

	//"development-of-internet-application/internal/app/dsn"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type Repository struct {
	Class *ClassRepository
	CalcRequest *CalcRequestRepository

}

func NewRepository() (*Repository, error) {

	dsnString := "host=localhost port=5432 user=postgres password=postgres dbname=dia_db sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		Class: NewClassRepository(db),
		CalcRequest: NewCalcRequestRepository(db),
	}, nil
}