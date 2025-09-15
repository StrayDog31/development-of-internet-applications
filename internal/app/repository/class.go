package repository

import (
	"fmt"
	"strings"
)

type ClassRepository struct {
}

func NewClassRepository() (*ClassRepository, error) {
	return &ClassRepository{}, nil
}

type SpectralClass struct {
	ID          int
	Name        string
	Temperature int
	Color       string
	Spectre     string
	Examples    string
	Image       string
}

var (
	spectralClasses = []SpectralClass{
		{
			ID:          1,
			Name:        "O",
			Temperature: 30000,
			Color:       "Голубой",
			Spectre:     "O5-O9",
			Examples:    "Горячие голубые звезды (Звезда Пистолет, 10 Ящерицы)",
			Image:       "/images/O-Class.png",
		},
		{
			ID:          2,
			Name:        "B",
			Temperature: 10000,
			Color:       "Бело-голубой",
			Spectre:     "B0-B9",
			Examples:    "Бело-голубые звезды (Ригель, Спика)",
			Image:       "/images/B-Class.png",
		},
		{
			ID:          3,
			Name:        "A",
			Temperature: 7500,
			Color:       "Белый",
			Spectre:     "A0-A9",
			Examples:    "Белые звезды (Сириус, Вега)",
			Image:       "/images/A-Class.png",
		},
		{
			ID:          4,
			Name:        "F",
			Temperature: 6000,
			Color:       "Желто-белый",
			Spectre:     "F0-F9",
			Examples:    "Желто-белые звезды (Процион, Канопус)",
			Image:       "/images/F-Class.png",
		},
		{
			ID:          5,
			Name:        "G",
			Temperature: 5200,
			Color:       "Желтый",
			Spectre:     "G0-G9",
			Examples:    "Желтые звезды (Солнце, Альфа Центавра A)",
			Image:       "/images/G-Class.png",
		},
		{
			ID:          6,
			Name:        "K",
			Temperature: 3700,
			Color:       "Оранжевый",
			Spectre:     "K0-K9",
			Examples:    "Оранжевые звезды (Арктур, Альдебаран)",
			Image:       "/images/K-Class.png",
		},
		{
			ID:          7,
			Name:        "M",
			Temperature: 2400,
			Color:       "Красный",
			Spectre:     "M0-M9",
			Examples:    "Красные карлики и гиганты (Бетельгейзе, Антарес)",
			Image:       "/images/M-Class.png",
		},
	}
)

func (*ClassRepository) GetClassByID(id int) (*SpectralClass, error) {
	if len(spectralClasses) == 0 {
		return nil, fmt.Errorf("масиив пуст")
	}

	for _, class := range spectralClasses {
		if class.ID == id {
			return &class, nil
		}
	}

	return nil, fmt.Errorf("не найдено")
}

func (*ClassRepository) GetClasses() ([]SpectralClass, error) {
	if len(spectralClasses) == 0 {
		return nil, fmt.Errorf("масиив пуст")
	}

	return spectralClasses, nil
}

func (repo *ClassRepository) GetClassByName(name string) ([]SpectralClass, error) {
	classes, err := repo.GetClasses()
	if err != nil {
		return []SpectralClass{}, err
	}

	var result []SpectralClass
	for _, class := range classes {
		if strings.Contains(strings.ToLower(class.Name), strings.ToLower(name)) {
			result = append(result, class)
		}
	}

	return result, nil
}