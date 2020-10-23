package responses

import (
	"time"
)

// Article response 数据结构
type Article struct {
	ID    int    `json:"id"`
	Title string `json:"Title"`
}

// ArticleRes response 数据结构
type ArticleRes struct {
	ID        int       `json:"id"`
	Summary   string    `json:"summary"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Category  Category  `json:"category"`
}

// ArticleRead 文章浏览次数
type ArticleRead struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Count int    `json:"count"`
}

// ArticleUserCount 用户发表文章次数
type ArticleUserCount struct {
	UserName string `json:"username"`
	Count    int    `json:"count"`
}
