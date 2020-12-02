package models

import (
	"community-blogger/internal/pkg/utils/constutil"
)

// User 定义用户表信息
type User struct {
	BaseModel
	Username string `gorm:"varchar(60) notnull 'username'" json:"username"`
	Password string `gorm:"varchar(60) notnull 'password'" json:"password"`
	Nickname string `form:"nickname" json:"nickname"`
	Mobile   string
	Token    string
}

// TableName  定义用户表名
func (User) TableName() string {
	return constutil.TablePrefix + "user"
}
