package requests

// FileInfo request 数据结构
type FileInfo struct {
	FileName string `form:"file_name" json:"file_name"`
}

// Download request 数据结构
type Download struct {
	FileURI string `form:"file_uri" json:"file_uri"`
}

// DeleteFile request 数据结构
type DeleteFile struct {
	FileURI string `form:"file_uri" json:"file_uri"`
}

// DeleteFiles request 数据结构
type DeleteFiles struct {
	FileURI []string `form:"file_uri" json:"file_uri"`
}
