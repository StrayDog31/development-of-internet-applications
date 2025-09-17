package repository

import (
	"development-of-internet-application/internal/app/ds"

	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{
		db: db,
	}
}

func (r *ClassRepository) GetClassByID(id uint64) (*ds.Class, error) {
	var class ds.Class
	err := r.db.Where("is_deleted = false").First(&class, id).Error
	if err != nil {
		return nil, err
	}

	return &class, nil
}

func (r *ClassRepository) GetClasses() ([]ds.Class, error) {
	var classes []ds.Class
	err := r.db.Where("is_deleted = false").Find(&classes).Error
	if err != nil {
		return nil, err
	}

	return classes, nil
}

func (r *ClassRepository) GetClassByName(name string) ([]ds.Class, error) {
	var classes []ds.Class
	err := r.db.
		Where("is_deleted = false AND name ILIKE ?", "%"+name+"%").
		Find(&classes).Error
	if err != nil {
		return nil, err
	}

	return classes, nil
}
