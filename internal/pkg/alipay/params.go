package alipay

import (
	"encoding/json"
)

// TradeParams trade参数
type TradeParams interface {
	// method
	Method() string
	// 针对特定需要加入到公共参数中的数据
	Params() map[string]string
	// 请求数据
	BizContent() string
}

// OriginBodySetter OriginBody
type OriginBodySetter interface {
	SetOriginBody(body []byte)
}

// marshal 转换json串
func marshal(obj interface{}) string {
	var bytes, err = json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bytes)
}
