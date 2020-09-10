package models

import (
	"time"
	"community-blogger/internal/pkg/utils/constutil"
)

// Integral 积分表
type Integral struct {
	ID        int `gorm:"column:id;primary_key" form:"id" json:"id"`
	UserID    int
	User      User      `gorm:"foreignKey:id;association_foreignkey:user_id"`
	Level     string    `db:"level" form:"level" json:"level"`
	DPoint    int64     `db:"dpoint" form:"dpoint" json:"dpoint"`
	DCoin     int64     `db:"dcoin" db:"dcoin" form:"dcoin" json:"dcoin"`
	CreatedAt time.Time `form:"created_at" json:"created_at" binding:"omitempty"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at" binding:"omitempty"`
}

// IntegralNote 积分记录表
type IntegralNote struct {
	ID        int `gorm:"column:id;primary_key" form:"id" json:"id"`
	UserID    int
	User      User      `gorm:"foreignKey:id;association_foreignkey:user_id"`
	Action    string    `db:"action" form:"action" json:"action"`
	DPoint    int64     `db:"dpoint" form:"dpoint" json:"dpoint"`
	DCoin     int64     `db:"dcoin" db:"dcoin" form:"dcoin" json:"dcoin"`
	CreatedAt time.Time `form:"created_at" json:"created_at" binding:"omitempty"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at" binding:"omitempty"`
}

// TableName 获取表名
func (Integral) TableName() string {
	return constutil.TablePrefix + "integral"
}

// TableName 获取表名
func (IntegralNote) TableName() string {
	return constutil.TablePrefix + "integral_note"
}
