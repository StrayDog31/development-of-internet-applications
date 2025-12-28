package ds

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"type:varchar(50);not null;unique"`
	PassHash string `gorm:"type:varchar(64);not null"`
	IsMod    bool   `gorm:"type:boolean;not null"`
}