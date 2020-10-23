package requests

import (
	"community-blogger/internal/pkg/responses"
	"time"
)

// Article request 数据结构
type Article struct {
	CategoryID int    `form:"category_id" json:"category_id"`
	Summary    string `db:"summary"`
	Title      string `db:"title"`
	UserName   string `form:"username" json:"username" binding:""`
}

// ArticleInfo request 数据结构
type ArticleInfo struct {
	ID int `form:"id" json:"id" binding:"required,min=1"`
}

// ArticleTop TOPN 结构
type ArticleTop struct {
	N int `form:"n" json:"n" binding:"required"`
}

// ArticleEdit request 数据结构
type ArticleEdit struct {
	ID         int    `form:"id" json:"id"`
	CategoryID int    `form:"category_id" json:"category_id"`
	Summary    string `db:"summary"`
	Title      string `db:"title"`
}

// ArticleCategoryEdit request 数据结构
type ArticleCategoryEdit struct {
	ID         int `form:"id" json:"id"`
	CategoryID int `form:"category_id" json:"category_id"`
}

// ArticleES request 数据结构
type ArticleES struct {
	ID        int                `json:"id"`
	Summary   string             `json:"summary"`
	Title     string             `json:"title"`
	Category  responses.Category `json:"category"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ArticleUserTop 用户发表文章数TON 结构
type ArticleUserTop struct {
	RankType string `form:"rank_type" json:"rank_type" binding:"required"`
	N        int    `form:"n" json:"n" binding:"required"`
}
