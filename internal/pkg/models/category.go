package models

import (
	"community-blogger/internal/pkg/utils/constutil"
	"time"
)

// Category 定义类别数据结构
type Category struct {
	ID        int       `gorm:"column:id;primary_key" form:"id" json:"id"`
	Name      string    `form:"name" json:"name"`
	Num       int       `form:"num" json:"num"`
	CreatedAt time.Time `form:"created_at" json:"created_at"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at"`
}

// TableName 获取类别表名
func (Category) TableName() string {
	return constutil.TablePrefix + "category"
}
