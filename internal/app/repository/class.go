package repository

import (
	"development-of-internet-applications/internal/app/ds"
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

type ClassFilter struct {
	Name    string
	Color   string
	Spectre string
	MinTemp int
	MaxTemp int
}

func (r *ClassRepository) GetClasses(filter *ClassFilter) ([]ds.Class, error) {
	var classes []ds.Class
	
	query := r.db.Where("is_deleted = false")
	
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Color != "" {
		query = query.Where("color ILIKE ?", "%"+filter.Color+"%")
	}
	if filter.Spectre != "" {
		query = query.Where("spectre ILIKE ?", "%"+filter.Spectre+"%")
	}
	if filter.MinTemp > 0 {
		query = query.Where("temperature >= ?", filter.MinTemp)
	}
	if filter.MaxTemp > 0 {
		query = query.Where("temperature <= ?", filter.MaxTemp)
	}
	
	err := query.Find(&classes).Error
	return classes, err
}

func (r *ClassRepository) GetClassByID(id uint64) (*ds.Class, error) {
	var class ds.Class
	err := r.db.Where("is_deleted = false").First(&class, id).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
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

func (r *ClassRepository) CreateClass(class *ds.Class) error {
	var maxID uint64
	err := r.db.Model(&ds.Class{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID).Error
	if err != nil {
		return err
	}

	class.ID = maxID + 1
	
	return r.db.Create(class).Error
}

func (r *ClassRepository) UpdateClass(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&ds.Class{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ClassRepository) DeleteClass(id uint64) error {
	return r.db.Model(&ds.Class{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ClassRepository) UpdateClassImage(id uint64, imagePath string) error {
	return r.db.Model(&ds.Class{}).Where("id = ?").Update("image", imagePath).Error
}