package repository

import (
	"development-of-internet-applications/internal/app/ds"
	"errors"
	"time"
	"gorm.io/gorm"
)

type MassRequestRepository struct {
	db *gorm.DB
}

func NewMassRequestRepository(db *gorm.DB) *MassRequestRepository {
	return &MassRequestRepository{
		db: db,
	}
}

type MassRequestFilter struct {
	Status    uint8
	StartDate time.Time
	EndDate   time.Time
}

func (r *MassRequestRepository) GetUserDraftRequest(userID uint64) (*ds.MassRequest, error) {
	var request ds.MassRequest
	err := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("user_id = ? AND status = ?", userID, 1).
		First(&request).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &request, err
}

func (r *MassRequestRepository) CreateRequest(request *ds.MassRequest) error {
	return r.db.Create(request).Error
}

func (r *MassRequestRepository) GetRequests(filter *MassRequestFilter) ([]ds.MassRequest, error) {
	var requests []ds.MassRequest
	
	query := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("status != ? AND status != ?", 1, 2)
	
	if filter.Status > 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if !filter.StartDate.IsZero() {
		query = query.Where("formed_at >= ?", filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query = query.Where("formed_at <= ?", filter.EndDate)
	}
	
	err := query.Order("created_at DESC").Find(&requests).Error
	return requests, err
}

func (r *MassRequestRepository) GetRequestByID(id uint64) (*ds.MassRequest, error) {
	var request ds.MassRequest
	err := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("id = ? AND status != ?", id, 2).
		First(&request).Error
	return &request, err
}

func (r *MassRequestRepository) UpdateRequest(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&ds.MassRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MassRequestRepository) DeleteRequest(id uint64) error {
	return r.db.Model(&ds.MassRequest{}).Where("id = ?", id).Update("status", 2).Error
}

func (r *MassRequestRepository) AddClassToRequest(requestID, classID uint64, luminosity, mass *uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingLink ds.MassRequestToClass
		err := tx.Where("request_id = ? AND class_id = ?", requestID, classID).First(&existingLink).Error
		if err == nil {
			updates := map[string]interface{}{}
			if luminosity != nil {
				updates["luminosity"] = luminosity
			}
			if mass != nil {
				updates["mass"] = mass
			}
			if len(updates) > 0 {
				return tx.Model(&ds.MassRequestToClass{}).
					Where("request_id = ? AND class_id = ?", requestID, classID).
					Updates(updates).Error
			}
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		item := ds.MassRequestToClass{
			RequestID:  requestID,
			ClassID:    classID,
			Luminosity: luminosity,
			Mass:       mass,
		}
		return tx.Create(&item).Error
	})
}

func (r *MassRequestRepository) RemoveClassFromRequest(requestID, classID uint64) error {
	return r.db.Where("request_id = ? AND class_id = ?", requestID, classID).
		Delete(&ds.MassRequestToClass{}).Error
}

func (r *MassRequestRepository) UpdateRequestClassItem(requestID, classID uint64, updates map[string]interface{}) error {
	return r.db.Model(&ds.MassRequestToClass{}).
		Where("request_id = ? AND class_id = ?", requestID, classID).
		Updates(updates).Error
}

func (r *MassRequestRepository) GetRequestItemsCount(requestID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&ds.MassRequestToClass{}).
		Where("request_id = ?", requestID).Count(&count).Error
	return count, err
}

func (r *MassRequestRepository) GetMassRequestByID(id uint64, userID uint64) (*ds.MassRequest, error) {
	var MassRequest ds.MassRequest
	err := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("id = ? AND user_id = ? AND status != 2", id, userID).
		First(&MassRequest).Error
	if err != nil {
		return nil, err
	}
	return &MassRequest, nil
}

func (r *MassRequestRepository) GetUserMassRequests(userID uint64) ([]ds.MassRequest, error) {
	var MassRequests []ds.MassRequest
	err := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("user_id = ? AND status != 2", userID).
		Find(&MassRequests).Error
	if err != nil {
		return nil, err
	}
	return MassRequests, nil
}

func (r *MassRequestRepository) CreateMassRequest(MassRequest *ds.MassRequest) error {
	return r.db.Create(MassRequest).Error
}

func (r *MassRequestRepository) AddClassToMassRequest(MassRequestID, classID uint64, luminosity, number *uint64) error {
	var existingLink ds.MassRequestToClass
	err := r.db.Where("request_id = ? AND class_id = ?", MassRequestID, classID).First(&existingLink).Error
	if err == nil {
		updates := map[string]interface{}{}
		if luminosity != nil {
			updates["luminosity"] = luminosity
		}
		if number != nil {
			updates["mass"] = number
		}
		if len(updates) > 0 {
			return r.db.Model(&ds.MassRequestToClass{}).
				Where("request_id = ? AND class_id = ?", MassRequestID, classID).
				Updates(updates).Error
		}
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	MassRequestToClass := ds.MassRequestToClass{
		RequestID:  MassRequestID,
		ClassID:    classID,
		Luminosity: luminosity,
		Mass:       number,
	}
	return r.db.Create(&MassRequestToClass).Error
}

func (r *MassRequestRepository) GetActiveMassRequestByUser(userID uint64) (*ds.MassRequest, error) {
	var MassRequest ds.MassRequest
	err := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("user_id = ? AND status = ?", userID, 1).
		Order("created_at DESC").
		First(&MassRequest).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &MassRequest, nil
}

func (r *MassRequestRepository) GetMassRequests(userID uint64, filter *MassRequestFilter) ([]ds.MassRequest, error) {
    var requests []ds.MassRequest
    
    query := r.db.
		Preload("MassRequestToClass").
		Preload("MassRequestToClass.Class").
		Where("user_id = ? AND status != ?", userID, 2)
    
    if filter.Status > 0 {
        query = query.Where("status = ?", filter.Status)
    }
    if !filter.StartDate.IsZero() {
        query = query.Where("formed_at >= ?", filter.StartDate)
    }
    if !filter.EndDate.IsZero() {
        query = query.Where("formed_at <= ?", filter.EndDate)
    }
    
    err := query.Order("created_at DESC").Find(&requests).Error
    return requests, err
}