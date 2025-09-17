package ds

type CalcRequestToClass struct {
	RequestID   uint64      `gorm:"primaryKey"`
	ClassID      uint64      `gorm:"primaryKey"`
	CalcRequest CalcRequest `gorm:"foreignKey:RequestID;references:ID"`
	Class        Class        `gorm:"foreignKey:ClassID;references:ID"`
	Lyminosity  *uint64      `gorm:"default:null"`
	Mass 		*uint64		`gorm:"default:null"`
	Number      *uint64      `gorm:"default:null"` 
}