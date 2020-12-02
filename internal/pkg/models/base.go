package models

import (
	"time"
)

const (
	FakeSize       = 9999
	DefaultPagSize = 15
)

// BaseModel .
type BaseModel struct {
	ID        uint      `gorm:"primary_key" json:"id" form:"id"`
	CreatedAt time.Time `form:"created_at" json:"created_at"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at"`
}
