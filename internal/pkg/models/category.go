package models

import (
	"community-blogger/internal/pkg/utils/constutil"
)

// Category 定义类别数据结构
type Category struct {
	BaseModel
	Name string `form:"name" json:"name"`
	Num  int    `form:"num" json:"num"`
}

// TableName 获取类别表名
func (Category) TableName() string {
	return constutil.TablePrefix + "category"
}
