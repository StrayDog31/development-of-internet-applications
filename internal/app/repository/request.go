package repository

import (
	"development-of-internet-application/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

type CalcRequestRepository struct {
	db *gorm.DB
}

func NewCalcRequestRepository(db *gorm.DB) *CalcRequestRepository {
	return &CalcRequestRepository{
		db: db,
	}
}

func (r *CalcRequestRepository) GetCalcRequestByID(id uint64, userID uint64) (*ds.CalcRequest, error) {
	var calcRequest ds.CalcRequest
	err := r.db.
		Preload("CalcRequestToClass").
		Preload("CalcRequestToClass.Class").
		Where("id = ? AND user_id = ? AND status != 2", id, userID).
		First(&calcRequest).Error
	if err != nil {
		return nil, err
	}

	return &calcRequest, nil
}

func (r *CalcRequestRepository) GetUserCalcRequests(userID uint64) ([]ds.CalcRequest, error) {
	var calcRequests []ds.CalcRequest
	err := r.db.
		Preload("CalcRequestToClass").
		Preload("CalcRequestToClass.Class").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&calcRequests).Error
	if err != nil {
		return nil, err
	}

	return calcRequests, nil
}

func (r *CalcRequestRepository) CreateCalcRequest(calcRequest *ds.CalcRequest) error {
	return r.db.Create(calcRequest).Error
}

func (r *CalcRequestRepository) AddClassToCalcRequest(calcRequestID, classID uint64, lyminosity, number *uint64) error {
	calcRequestToClass := ds.CalcRequestToClass{
		RequestID:  calcRequestID,
		ClassID:    classID,
		Lyminosity: lyminosity,
		Number:     number,
	}
	return r.db.Create(&calcRequestToClass).Error
}

func (r *CalcRequestRepository) DeleteRequest(requestID uint64, userID uint64) error {
	result := r.db.Exec(`
		UPDATE calc_requests 
		SET status = 2, closed_at = NOW() 
		WHERE id = ? AND user_id = ? AND status = 1
	`, requestID, userID)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
func (r *CalcRequestRepository) GetActiveCalcRequestByUser(userID uint64) (*ds.CalcRequest, error) {
	var calcRequest ds.CalcRequest
	err := r.db.
		Preload("CalcRequestToClass").
		Where("user_id = ? AND status = ?", userID, 1).
		Order("created_at DESC").
		First(&calcRequest).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &calcRequest, nil
}

func (r *CalcRequestRepository) AddClassToRequest (classId uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lamp ds.Class
		err := r.db.First(&lamp, classId).Error
		if err != nil {
			return err
		}

		var calcRequest ds.CalcRequest
		err = r.db.
			Where("status = 1 AND user_id = ?", userID).
			Take(&calcRequest).Error
		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			calcRequest = ds.CalcRequest{
				User: ds.User{ID: userID},
			}
			err := r.db.Create(&calcRequest).Error
			if err != nil {
				return err
			}
		}

		calcRequestToLamp := ds.CalcRequestToClass{
			RequestID: calcRequest.ID,
			ClassID:    classId,
		}
		r.db.Create(&calcRequestToLamp)

		return nil
	})
}