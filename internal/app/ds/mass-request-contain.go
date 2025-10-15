package ds

type MassRequestToClass struct {
	RequestID   uint64      `gorm:"primaryKey"`
	ClassID      uint64      `gorm:"primaryKey"`
	CalcRequest MassRequest `gorm:"foreignKey:RequestID;references:ID"`
	Class        Class        `gorm:"foreignKey:ClassID;references:ID"`
	Luminosity  *uint64      `gorm:"default:null"`
	Mass 		*uint64		`gorm:"default:null"`
}