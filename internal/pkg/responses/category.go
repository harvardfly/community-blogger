package responses

import (
	"time"
)

// Category response 数据结构
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Num       int       `json:"num"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
