package dbutil

import (
	"community-blogger/internal/pkg/utils/constutil"
	"github.com/jinzhu/gorm"
	"strconv"
)

// ScopeFunc 定义Scope方法
type ScopeFunc func(db *gorm.DB) *gorm.DB

// PageOrder 返回gorm分页方法
func PageOrder(p map[string]int, o map[string]string) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		limit, ok := p["limit"]
		if !ok {
			limit = constutil.DefaultPageSize
		}
		if limit > constutil.MaxPageSize {
			limit = constutil.MaxPageSize
		}
		offset, ok := p["offset"]
		if !ok {
			offset = 0
		}
		if offset < 0 {
			offset = 0
		}
		str := ""
		field, ok := o["order"]
		if !ok || field == "" {
			return db.Offset(offset * limit).Limit(limit).Order("id desc")
		}
		b, ok := o["desc"]
		if !ok || b == "1" {
			str = field + " " + "desc"
		} else {
			str = field + " " + "asc"
		}
		return db.Offset(offset * limit).Limit(limit).Order(str)
	}
}

// PageOptimize 分页优化 优化原limit offset --> 先limit offset取id 再查询返回 10w+数据时优化10倍以上
func PageOptimize(p map[string]int, o map[string]string, tb string, sqlWhere string) string {
	limit, ok := p["limit"]
	if !ok {
		limit = constutil.DefaultPageSize
	}
	if limit > constutil.MaxPageSize {
		limit = constutil.MaxPageSize
	}
	offset, ok := p["offset"]
	if !ok {
		offset = 0
	}
	if offset < 0 {
		offset = 0
	}
	order := ""
	field, ok := o["order"]
	if !ok || field == "" {
		order = "id desc"
	} else {
		order = field + " " + "asc"
		b, ok := o["desc"]
		if !ok || b == "1" {
			order = field + " " + "desc"
		}
	}
	offset = offset * limit
	sqlStr := "SELECT a.* FROM " + tb + " a, (select id from " + tb + sqlWhere + " limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset) + " ) b where a.id=b.id ORDER BY " + order
	return sqlStr
}
