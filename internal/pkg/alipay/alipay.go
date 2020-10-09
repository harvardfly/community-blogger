package alipay

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AliPay 支付宝支付 struct
type AliPay struct {
	appId           string
	apiDomain       string
	publicKey       []byte
	privateKey      []byte
	AliPayPublicKey []byte
	client          *http.Client
	SignType        string
}

// NewAliPay 初始化AliPay
func NewAliPay(appId string, publicKey, privateKey []byte, isProduction bool) (client *AliPay) {
	client = &AliPay{
		appId:           appId,
		privateKey:      privateKey,
		publicKey:       publicKey,
		AliPayPublicKey: publicKey,
		client:          http.DefaultClient,
		SignType:        SignTypeRsa2,
	}

	apiDomain := SandboxApiURL
	if isProduction {
		apiDomain = ProductionApiURL
	}
	client.apiDomain = apiDomain
	return client
}

// URLValues 支付链接信息
func (a *AliPay) URLValues(param TradeParams) (value url.Values, err error) {
	var p = url.Values{}
	p.Add("app_id", a.appId)
	p.Add("method", param.Method())
	p.Add("format", FORMAT)
	p.Add("charset", CHARSET)
	p.Add("sign_type", a.SignType)
	p.Add("timestamp", time.Now().Format(TimeFormat))
	p.Add("version", VERSION)

	if len(param.BizContent()) > 0 {
		p.Add("biz_content", param.BizContent())
	}

	var ps = param.Params()
	if ps != nil {
		for key, value := range ps {
			p.Add(key, value)
		}
	}

	var keys = make([]string, 0, 0)
	for key := range p {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var sign string
	if a.SignType == SignTypeRsa {
		sign, err = signRSA(keys, p, a.privateKey)
	} else {
		sign, err = signRSA2(keys, p, a.privateKey)
	}
	p.Add("sign", sign)

	if err != nil {
		return nil, err
	}
	return p, nil
}

// doRequestWithVerify 校验
func (a *AliPay) doRequestWithVerify(method string, param TradeParams, results interface{}) error {
	body, err := a.doRequest(method, param, results)
	if err != nil {
		return err
	}

	if len(a.AliPayPublicKey) > 0 {
		var rootNodeName = strings.Replace(param.Method(), ".", "_", -1) + ResponseSuffix
		content, sign := parseJSONSource(string(body), rootNodeName)
		if ok, err := verifyResponseData([]byte(content), a.SignType, sign, a.AliPayPublicKey); ok == false {
			return err
		}
	}
	return nil
}

// doRequest 执行请求
func (a *AliPay) doRequest(method string, param TradeParams, results interface{}) (body []byte, err error) {
	var buf io.Reader
	if param != nil {
		p, err := a.URLValues(param)
		if err != nil {
			return nil, err
		}
		buf = strings.NewReader(p.Encode())
	}

	req, err := http.NewRequest(method, a.apiDomain, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if results != nil {
		err = json.Unmarshal(body, results)
		if err != nil {
			return
		}
	}

	if t, ok := results.(OriginBodySetter); ok {
		t.SetOriginBody(body)
	}

	return body, nil
}

// parseJSONSource 解析json
func parseJSONSource(body, nodeName string) (content, sign string) {
	body = strings.TrimPrefix(body, "{")
	body = strings.TrimSuffix(body, "}")

	contentStartIndex := strings.Index(body, "{")
	contentEndIndex := strings.LastIndex(body, "}") + 1
	if contentEndIndex <= contentStartIndex {
		return
	}
	content = body[contentStartIndex:contentEndIndex]

	sep := fmt.Sprintf("\"%s\":", nodeName) + content
	splits := strings.Split(body, sep)

	for _, split := range splits {
		if strings.Contains(split, "sign") {
			body = split
			break
		}
	}

	body = strings.TrimSuffix(body, `"`)
	body = strings.TrimSuffix(body, `",`)

	body = strings.TrimPrefix(body, `,"sign":"`)
	body = strings.TrimPrefix(body, `"sign":"`)
	sign = body
	return
}

// signRSA2 RSA2签名
func signRSA2(keys []string, param url.Values, privateKey []byte) (s string, err error) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}

	var src = strings.Join(pList, "&")
	sig, err := SignPKCS1v15([]byte(src), privateKey, crypto.SHA256)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

// signRSA RSA签名
func signRSA(keys []string, param url.Values, privateKey []byte) (s string, err error) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}

	var src = strings.Join(pList, "&")
	sig, err := SignPKCS1v15([]byte(src), privateKey, crypto.SHA1)
	if err != nil {
		return "", err
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s, nil
}

// verifySign 校验签名
func verifySign(req *http.Request, key []byte) (ok bool, err error) {
	sign, err := base64.StdEncoding.DecodeString(req.Form.Get("sign"))
	signType := req.Form.Get("sign_type")
	if err != nil {
		return false, err
	}

	var keys = make([]string, 0, 0)
	for key, value := range req.Form {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = req.Form.Get(key)
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")

	if signType == SignTypeRsa {
		err = VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA1)
	} else {
		err = VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA256)
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

// verifyResponseData 校验返回数据
func verifyResponseData(data []byte, signType, sign string, key []byte) (ok bool, err error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	if signType == SignTypeRsa {
		err = VerifyPKCS1v15(data, signBytes, key, crypto.SHA1)
	} else {
		err = VerifyPKCS1v15(data, signBytes, key, crypto.SHA256)
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
