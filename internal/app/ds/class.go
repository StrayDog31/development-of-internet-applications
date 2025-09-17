package ds

type Class struct {
	ID                 uint64    `gorm:"primaryKey"`
	Name              string     `gorm:"type:varchar(100);not null"`
	Temperature int				 `gorm:"not null"`
	Color       string			`gorm:"type:varchar(100);not null"`
	Spectre     string			`gorm:"type:varchar(100);not null"`
	Examples    string			`gorm:"type:text;not null"`
	IsDeleted       bool     `gorm:"boolean;not null"`
	Image           string      `gorm:"type:varchar(100)"`
	CalcRequestToClass  []CalcRequestToClass `gorm:"foreignKey:ClassID"`
}