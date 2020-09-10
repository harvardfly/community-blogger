package httputil

// Error 定义异常返回
func Error(err error, msg string) interface{} {
	return struct {
		Err string `json:"err"`
		Msg string `json:"msg"`
	}{
		Err: err.Error(),
		Msg: msg,
	}
}

// Success 定义成功返回
func Success(data interface{}) interface{} {
	return struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}
}

// ReturnJSON 定义返回json
func ReturnJSON(code int64, message string, data interface{}) interface{} {
	return struct {
		Code    int64       `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
