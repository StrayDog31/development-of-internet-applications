package repository

import (
	"fmt"
)

type CalcRequestRepository struct {
}

func NewCalcRequestRepository() (*CalcRequestRepository, error) {
	return &CalcRequestRepository{}, nil
}

type CalcRequest struct {
	ID    int
	Stars []DefinedStar
}

type DefinedStar struct {
	Class       SpectralClass
	Luminosity  float64 
	Mass        float64 
	Temperature int
	Type        string
}

var (
	calcRequests = []CalcRequest{
		{
			ID: 1,
			Stars: []DefinedStar{
				{
					Class:       SpectralClass{ID: 1, Name: "O", Temperature: 30000, Color: "Голубой"},
					Luminosity:  250,
					Mass:        33,
					Temperature: 35500,
					Type:        "Сверхгигант",
				},
				{
					Class:       SpectralClass{ID: 2, Name: "B", Temperature: 10000, Color: "Бело-голубой"},
					Luminosity:  120,
					Mass:        18,
					Temperature: 15000,
					Type:        "Гигант",
				},
			},
		},
	}
)

func (*CalcRequestRepository) GetCalcRequestViewByID(id int, classRepo *ClassRepository) (*CalcRequest, error) {
	if len(calcRequests) == 0 {
		return nil, fmt.Errorf("массив пуст")
	}

	for _, req := range calcRequests {
		if req.ID == id {
			calcRequest := req
			
			for i := range calcRequest.Stars {
				if classRepo != nil {
					fullClass, err := classRepo.GetClassByID(calcRequest.Stars[i].Class.ID)
					if err == nil {
						calcRequest.Stars[i].Class = *fullClass
					}
				}
			}
			
			return &calcRequest, nil
		}
	}
	
	return nil, fmt.Errorf("не найдено")
}