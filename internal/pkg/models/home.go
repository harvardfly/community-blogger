package models

import (
	"time"
	"community-blogger/internal/pkg/utils/constutil"
)

// Home 定义home数据结构
type Home struct {
	ID          int
	URL         string
	Img         string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName 获取home表名
func (Home) TableName() string {
	return constutil.TablePrefix + "home"
}

/*
CREATE TABLE `blog_home` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
	`url` varchar(255) NOT NULL COMMENT '跳转地址',
	`img` varchar(255) NOT NULL COMMENT '图片地址',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `description` text COMMENT '描述',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=480 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='博客首页';
*/
