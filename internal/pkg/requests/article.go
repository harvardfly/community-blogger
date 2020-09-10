package requests

// Article request 数据结构
type Article struct {
	CategoryID int    `form:"category_id" json:"category_id"`
	Summary    string `db:"summary"`
	Title      string `db:"title"`
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
