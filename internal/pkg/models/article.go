package models

import (
	"community-blogger/internal/pkg/utils/constutil"
	"time"
)

// Article 定义文章数据结构
type Article struct {
	ID         int `gorm:"column:id;primary_key" form:"id" json:"id"`
	CategoryID int
	Category   Category  `gorm:"foreignKey:id;association_foreignkey:category_id"`
	Summary    string    `db:"summary"`
	Title      string    `db:"title"`
	CreatedAt  time.Time `form:"created_at" json:"created_at" binding:"omitempty"`
	UpdatedAt  time.Time `form:"updated_at" json:"updated_at" binding:"omitempty"`
}

// TableName 获取表名
func (Article) TableName() string {
	return constutil.TablePrefix + "article"
}
