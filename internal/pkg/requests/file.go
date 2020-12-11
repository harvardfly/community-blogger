package requests

// Home request 数据结构
type FileInfo struct {
	FileName string `form:"file_name" json:"file_name"`
}
