package responses

import (
	"time"
)

// Home response 数据结构
type Home struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Img         string    `json:"img"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
