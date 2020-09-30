package ginalipay

import (
	"bufio"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
)

// TradeNotification https://docs.open.alipay.com/203/105286
type TradeNotification struct {
	NotifyTime        string `json:"notify_time"`         // 通知时间,通知的发送时间。格式为yyyy-MM-dd HH:mm:ss
	NotifyType        string `json:"notify_type"`         // 通知类型
	NotifyId          string `json:"notify_id"`           // 通知校验ID
	AppId             string `json:"app_id"`              // 开发者的app_id
	Charset           string `json:"charset"`             // 编码格式
	Version           string `json:"version"`             // 接口版本
	SignType          string `json:"sign_type"`           // 签名类型
	Sign              string `json:"sign"`                // 签名
	TradeNo           string `json:"trade_no"`            // 支付宝交易号
	OutTradeNo        string `json:"out_trade_no"`        // 商户订单号
	OutBizNo          string `json:"out_biz_no"`          // 商户业务号
	BuyerId           string `json:"buyer_id"`            // 买家支付宝用户号
	BuyerLogonId      string `json:"buyer_logon_id"`      // 买家支付宝账号
	SellerId          string `json:"seller_id"`           // 卖家支付宝用户号
	SellerEmail       string `json:"seller_email"`        // 卖家支付宝账号
	TradeStatus       string `json:"trade_status"`        // 交易状态
	TotalAmount       string `json:"total_amount"`        // 订单金额
	ReceiptAmount     string `json:"receipt_amount"`      // 实收金额
	InvoiceAmount     string `json:"invoice_amount"`      // 开票金额
	BuyerPayAmount    string `json:"buyer_pay_amount"`    // 付款金额
	PointAmount       string `json:"point_amount"`        // 集分宝金额
	RefundFee         string `json:"refund_fee"`          // 总退款金额
	Subject           string `json:"subject"`             // 总退款金额
	Body              string `json:"body"`                // 商品描述
	GmtCreate         string `json:"gmt_create"`          // 交易创建时间
	GmtPayment        string `json:"gmt_payment"`         // 交易付款时间
	GmtRefund         string `json:"gmt_refund"`          // 交易退款时间
	GmtClose          string `json:"gmt_close"`           // 交易结束时间
	FundBillList      string `json:"fund_bill_list"`      // 支付金额信息
	PassbackParams    string `json:"passback_params"`     // 回传参数
	VoucherDetailList string `json:"voucher_detail_list"` // 优惠券信息
}

// GetTradeNotificationByBody 获取交易通知 body测试
func (a *AliPay) GetTradeNotificationByBody(body string) (*TradeNotification, error) {
	fakeURL := "http://fake.alipay.notify/?" + body
	method := http.MethodPost
	_, err := http.ReadRequest(bufio.NewReader(strings.NewReader(method + " " + fakeURL + " HTTP/1.0\r\n\r\n")))
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(http.MethodPost, fakeURL, strings.NewReader(""))
	if req == nil {
		return nil, errors.New("make server request failed")
	}
	return a.GetTradeNotification(req)
}

// GetTradeNotification 获取交易通知
func (a *AliPay) GetTradeNotification(req *http.Request) (*TradeNotification, error) {
	if req == nil {
		return nil, errors.New("invalid req")
	}

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	noti := &TradeNotification{
		NotifyTime:        req.FormValue("notify_time"),
		NotifyType:        req.FormValue("notify_type"),
		NotifyId:          req.FormValue("notify_id"),
		AppId:             req.FormValue("app_id"),
		Charset:           req.FormValue("charset"),
		Version:           req.FormValue("version"),
		SignType:          req.FormValue("sign_type"),
		Sign:              req.FormValue("sign"),
		TradeNo:           req.FormValue("trade_no"),
		OutTradeNo:        req.FormValue("out_trade_no"),
		OutBizNo:          req.FormValue("out_biz_no"),
		BuyerId:           req.FormValue("buyer_id"),
		BuyerLogonId:      req.FormValue("buyer_logon_id"),
		SellerId:          req.FormValue("seller_id"),
		SellerEmail:       req.FormValue("seller_email"),
		TradeStatus:       req.FormValue("trade_status"),
		TotalAmount:       req.FormValue("total_amount"),
		ReceiptAmount:     req.FormValue("receipt_amount"),
		InvoiceAmount:     req.FormValue("invoice_amount"),
		BuyerPayAmount:    req.FormValue("buyer_pay_amount"),
		PointAmount:       req.FormValue("point_amount"),
		RefundFee:         req.FormValue("refund_fee"),
		Subject:           req.FormValue("subject"),
		Body:              req.FormValue("body"),
		GmtCreate:         req.FormValue("gmt_create"),
		GmtPayment:        req.FormValue("gmt_payment"),
		GmtRefund:         req.FormValue("gmt_refund"),
		GmtClose:          req.FormValue("gmt_close"),
		FundBillList:      req.FormValue("fund_bill_list"),
		PassbackParams:    req.FormValue("passback_params"),
		VoucherDetailList: req.FormValue("voucher_detail_list"),
	}

	if len(noti.NotifyId) == 0 {
		return nil, errors.New("invalid notification")
	}

	ok, err := verifySign(req, a.AliPayPublicKey)
	if ok {
		return noti, nil
	}
	return nil, err
}
