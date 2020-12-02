package models

import (
	"community-blogger/internal/pkg/utils/constutil"
)

// Article 定义文章数据结构
type Article struct {
	BaseModel
	CategoryID int
	Category   Category `gorm:"foreignKey:id;association_foreignkey:category_id"`
	Summary    string   `db:"summary"`
	Title      string   `db:"title"`
}

// TableName 获取表名
func (Article) TableName() string {
	return constutil.TablePrefix + "article"
}
