package requests

// HomeList request 数据结构
type HomeList struct {
	ID          int               `form:"id" json:"id"`
	Limit       int               `form:"limit" json:"limit" binding:"omitempty"`
	Page        int               `form:"page" json:"page" binding:"omitempty"`
	Order       string            `form:"order" json:"order" binding:"omitempty"`
	Description string            `form:"description" json:"description" binding:"omitempty"`
	Title       string            `form:"title" json:"title" binding:"omitempty"`
	Paginator   map[string]int    `binding:"-"`
	Orders      map[string]string `binding:"-"`
}

// Home request 数据结构
type Home struct {
	URL         string `form:"url" json:"url" binding:"omitempty"`
	Img         string `form:"img" json:"img" binding:"required"`
	Title       string `form:"title" json:"title" binding:"required"`
	Description string `form:"description" json:"description" binding:"omitempty"`
}
