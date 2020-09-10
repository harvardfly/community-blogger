package responses

// UserInfo 定义用户返回结构体
type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// LoginResponse 定义登录返回结构体
type LoginResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"accessToken"`
	ExpireAt    int64  `json:"expireAt"`
	TimeStamp   int64  `json:"timeStamp"`
}

// RegisterResponse 定义注册返回结构体
type RegisterResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
