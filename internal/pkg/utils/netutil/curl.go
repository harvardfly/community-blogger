package netutil

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Curl 定义方法执行curl
func Curl(method, u string, params io.Reader, logger *zap.Logger) (*http.Response, error) {
	client := http.Client{
		Timeout: 30,
	}
	var (
		req *http.Request
		err error
	)

	if strings.ToLower(method) == "post" {
		req, err = http.NewRequest("POST", u, params)
	} else {
		p, err := ioutil.ReadAll(params)
		if err != nil {
			return nil, err
		}
		path, err := url.ParseQuery(string(p))
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("GET", strings.Trim(u, "/")+"?"+path.Encode(), nil)
	}
	if err != nil {
		logger.Error("http newRequest失败", zap.Error(err))
		return nil, err
	}
	return client.Do(req)
}

// HTTPCurl 定义以http的方式执行curl
func HTTPCurl(method, u string, data io.Reader, header map[string]string, t time.Duration, logger *zap.Logger) (*http.Response, error) {
	client := http.Client{Timeout: t}
	var (
		req *http.Request
		err error
	)
	switch strings.ToLower(method) {
	case "post", "put", "delete", "head":
		req, err = http.NewRequest(method, u, data)
	case "get":
		if data != nil {
			p, err := ioutil.ReadAll(data)
			if err != nil {
				logger.Error("读取数据失败", zap.Error(err))
				return nil, err
			}
			path, err := url.ParseQuery(string(p))
			if err != nil {
				logger.Error("数据转换url path失败", zap.Error(err))
				return nil, err
			}
			if strings.Contains(u, "?") {
				u += "&" + path.Encode()
			} else {
				u += "?" + path.Encode()
			}
		}
		req, err = http.NewRequest("GET", u, nil)
	default:
		return nil, errors.New("不支持该请求方式")
	}
	if err != nil {
		logger.Error("http newRequest失败", zap.Error(err))
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return client.Do(req)
}
