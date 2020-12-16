package responses

// FileInfo response 数据结构
type FileInfo struct {
	Hash     string `json:"hash"`
	Fsize    string `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	URI      string `json:"uri"`
}
