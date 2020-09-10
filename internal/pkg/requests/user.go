package requests

// UserCreate 创建用户的请求结构体
type UserCreate struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"limit" binding:"required"`
	Password2 string `form:"password2" json:"limit" binding:"required"`
	Nickname  string `form:"nick_name" json:"nick_name" binding:"omitempty"`
	Mobile    string `form:"mobile" json:"mobile" binding:"required"`
	Token     string `form:"token" json:"token"`
}

// User 用户请求结构体
type User struct {
	ID int `form:"id" json:"id"`
}

// UserToken 获取用户Token请求结构体
type UserToken struct {
	Token string `form:"token" json:"token"`
}

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `form:"nick_name" json:"nick_name" binding:"omitempty"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required"`
}
