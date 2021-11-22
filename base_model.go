package db

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"created_at" json:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updated_at" form:"updated_at" query:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `gorm:"primaryKey" json:"id"`
	IsActive  bool           `gorm:"is_active" json:"is_active" form:"is_active" query:"is_active"`
}
