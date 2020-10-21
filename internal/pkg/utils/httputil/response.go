package httputil

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

// 定义一个全局翻译器T
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// Error 定义异常返回
func Error(err error, msg string) interface{} {
	if err := InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		return struct {
			Err string `json:"err"`
			Msg string `json:"msg"`
		}{
			Err: err.Error(),
			Msg: msg,
		}
	} else {
		// validator.ValidationErrors类型错误则进行翻译
		return struct {
			Err map[string]string `json:"err"`
			Msg string            `json:"msg"`
		}{
			Err: removeTopStruct(errs.Translate(trans)),
			Msg: msg,
		}
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

// GetReferDomain 获取RequestURI或Referer的域名
func GetReferDomain(requestURI string) string {
	tmpURI := requestURI
	if find := strings.Contains(requestURI, "//"); find {
		tmpURI = strings.Split(requestURI, "//")[1]
	}
	res := strings.Split(tmpURI, "/")[0]
	return res
}
