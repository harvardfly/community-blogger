package alipay

import (
	"context"
	"os"
	"testing"
	"time"

	"gotest.tools/assert"
)

var (
	isSandBox = true
	// sandbox
	appID     = "2016092300580717"
	publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp3TKHLE439G6oVyjq4v0nnFmLoONphBwd6rulG1dTKc5qybDJykf/dPqJnULtPUs1oRlgmoBK0NCUB2OPuYpM2v2Cfvvm4zsXW9DuzYaItFB+Nd5KmZU7hMK6oNHOgCKiNW1pYILOwoJaKurRn/UbffVDzxeEiBoLLk5kpszbu53LngLYG7RQhGkjh6/cWuy3o4MdXuee66v1tDOEK2U/5lOx9r+Y4SY7b6pT8XW95/KFROze8Wn8mXbrJWfDZ6Z8Ee0ny5miKXthhdrE/KqNsHRBKTS/8aVdmAedLKw7lAoC28jZssKN0YsHHcWd9++jVPLmPK7XgQmpd73a1B0LQIDAQAB
-----END PUBLIC KEY-----`)

	privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAp3TKHLE439G6oVyjq4v0nnFmLoONphBwd6rulG1dTKc5qybDJykf/dPqJnULtPUs1oRlgmoBK0NCUB2OPuYpM2v2Cfvvm4zsXW9DuzYaItFB+Nd5KmZU7hMK6oNHOgCKiNW1pYILOwoJaKurRn/UbffVDzxeEiBoLLk5kpszbu53LngLYG7RQhGkjh6/cWuy3o4MdXuee66v1tDOEK2U/5lOx9r+Y4SY7b6pT8XW95/KFROze8Wn8mXbrJWfDZ6Z8Ee0ny5miKXthhdrE/KqNsHRBKTS/8aVdmAedLKw7lAoC28jZssKN0YsHHcWd9++jVPLmPK7XgQmpd73a1B0LQIDAQABAoIBAEjU6vL/wZTXSyzTdfwuqv4epCqm3PzVOZVSquGzj1i/gr2F5msp39guSzDex3C1EgNbIitOn6OJZVYjBLMmt9S9qA0/nj8xU7xvoC3Uohlymhb44KIoT9gcQCsvXdNEWPyatp59qRTMkLsNrzjCcEpD+E7gGqoXnjeVeAzjo8MjD6O1/jMH4Diug7Ulopm+avWR8c4gAlBopcapDJwgLT5j8t9hlWxyRyoYJif+0iwGmzEtWfar8/j4t8IJP9XGOI891ZbRNOhGPy58tVN4ET4Dq6hEtgqrQFFUTweomGcBjK84L85UTncZrJ/nV4UDcQmirVLEfuvxtT33B6q3AQECgYEA4w9gxrqsQkgG+bETbZNLtEYDGmL2g49uoa0bf4VkDUr3kE4tFfQVVL26Z0PG9ASOdNvxWXFtfGGQPSfBOeWKVRxAcGN50TUPn1NyGBVFo+TaiTpjQePmyivuhStYjpM4ER72kv71Dfq9d1TXAH8Yd7GtkwhFMr92OaVgUNqMT5UCgYEAvMyg/JK8taB05L3HEKJD4IblRB+5Vfq8uwoIeq9icxj+tL7m/DYsV4w0rq6eTR4ZPk4cADE71nOfFC8IQPQBdZCl3fMTiPSAd9P8J/SRiSzbVqCm49A0eljYDZNJ5UmLlysUfIRAWfg0LXkO4g6qe4zMBe8kye05WSNydR4IzDkCgYEAhbdxs7cNaP9H9FX+7gHXjHPKsghjynh2m4n5brDciiODBZ20WYBj64LMOrIkgWIJjvJUAOuWobBHJGy8E1+FhrfbjxRWEglgiOC4iUxFtc45kKUs/Qm8yYTzs8MiJNy4IQUOCuVQ5YmreJIjB+zROPQPF07AibFNa4dj4FnEVb0CgYEAnsu18/ovgsxsxR/h2NnCIY0UNJJCPlDggKjVrOrq8Ufqo3eVrDicXx5sCSiRuOdB8CeeYYMHgz5IZJ+SX59bwthgyidzHNQZYbAI/Eo2RpxV96yz8hTirq0fO3vQwWt0Xzc6yegqgAHFUSHrJLOVLMmlqLAVz0kZ3SF1WZBjcfECgYEAyuU8vfl3nNcG71DCM+HvCu/mePrPMnie08dzPcUS7oC7LNm0ZKB9MuHM1LKfCKoifbhRfUlOfTMyi8x8t8uVVSrWZn7PutKI+k86XR7B71Bs6S1DXUjI8562vlzA9EJq+Umlos6eaR3tZNXiUgekFtb2kJ7+81Ab1TTTqR+IoVM=
-----END RSA PRIVATE KEY-----`)
	appIDDev      = ""
	publicKeyDev  = []byte("")
	privateKeyDev = []byte("")
)

var client *AliPay

func TestMain(m *testing.M) {
	isProduction := false
	if !isSandBox {
		isProduction = true
		appID = appIDDev
		publicKey = publicKeyDev
		privateKey = privateKeyDev
	}
	client = NewAliPay(appID, publicKey, privateKey, isProduction)
	os.Exit(m.Run())
}

func TestTradeAppPay(t *testing.T) {
	ret, err := client.TradeAppPay(context.Background(), &TradeAppPayRequest{
		TotalAmount: "1",
		Subject:     "测试商品",
		OutTradeNo:  "00000000000000001",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

// done
func TestTradeWapPay(t *testing.T) {
	ret, err := client.TradeWapPay(context.Background(), &TradeWapPayRequest{
		TotalAmount:    "100",
		Subject:        "测试商品",
		OutTradeNo:     "00906",
		ProductCode:    "QUICK_WAP_WAY",
		QuitURL:        "http://www.baidu.com",
		NotifyURL:      "",
		TimeoutExpress: "30m",
		GoodsType:      "1",
		TimeExpire:     time.Now().Add(time.Minute * 30).Format("2006-01-02 15:04"),
	})
	assert.Equal(t, nil, err)
	t.Logf("WapPayTargetURL:%s", ret.TargetURL)
}

// done
func TestTradePagePay(t *testing.T) {
	ret, err := client.TradePagePay(context.Background(), &TradePagePayRequest{
		TotalAmount:    "100",
		Subject:        "测试商品",
		OutTradeNo:     "00099",
		ProductCode:    "FAST_INSTANT_TRADE_PAY",
		NotifyURL:      "",
		TimeoutExpress: "30m",
		GoodsType:      "1",
		TimeExpire:     time.Now().Add(time.Minute * 30).Format("2006-01-02 15:04"),
	})
	assert.Equal(t, nil, err)
	t.Logf("PagePayTargetURL:%s", ret.TargetURL)
}

// done
func TestTradeQuery(t *testing.T) {
	ret, err := client.TradeQuery(context.Background(), &TradeQueryRequest{
		OutTradeNo: "00906",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

func TestTradeClose(t *testing.T) {
	ret, err := client.TradeClose(context.Background(), &TradeCloseRequest{
		OutTradeNo: "004",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

// done
func TestTradeCancel(t *testing.T) {
	ret, err := client.TradeCancel(context.Background(), &TradeCancelRequest{
		OutTradeNo: "004",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

// done
func TestTradeRefund(t *testing.T) {
	ret, err := client.TradeRefund(context.Background(), &TradeRefundRequest{
		OutTradeNo:   "005",
		RefundAmount: 50,
		OutRequestNo: "request_005_01",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

// done
func TestTradeRefundQuery(t *testing.T) {
	ret, err := client.TradeRefundQuery(context.Background(), &TradeRefundQueryRequest{
		OutTradeNo:   "005",
		OutRequestNo: "request_003_03",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

// done
func TestBillDownloadURLQuery(t *testing.T) {
	ret, err := client.BillDownloadURLQuery(context.Background(), &BillDownloadRequest{
		BillType: "trade",
		BillDate: "2019-10-21",
	})
	assert.Equal(t, nil, err)
	t.Logf("%#v", *ret)
}

func TestParseJSONSource(t *testing.T) {
	body := `{"sign":"UpBhxKyuxGEMRethWKXj13Vj3rG7YJb+tC8YQdISwwerv7Hnr87O/hpSpEOMj30IrVfrUs4IMaaHHI+VaNM9WlmdqkJMXcqrVkh0VTzLykYWsfI08EbYMkbj0MyNFO27JjbALq7kknFzjj2p3lf6X+zo1nE6Wh678yFZ3esVUIkI7YxB0sQQvxkG/kMefwDhGwigtljXf6QEylzKGHDQy9hlCOUHFDSM5XFd1DJQuFqhttG/ujOkRTSJy6DaV/kjjxaA0iLhDoCqcGHBxTN9WQtYV3I+hIHhOLUtmLmQzHTypA7bOOgKGv6zL+3zObftEIPdfN/pii2D5k6vbtmSbg==","alipay_trade_query_response":{"code":"10000","msg":"Success","buyer_logon_id":"hfj***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088102179843724","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"003","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2019-10-22 15:59:49","total_amount":"100.00","trade_no":"2019102222001443721000074028","trade_status":"TRADE_CLOSED"}}`
	rootNodeName := "alipay_trade_query_response"
	content, sign := parseJSONSource(body, rootNodeName)
	assert.Equal(t, `{"code":"10000","msg":"Success","buyer_logon_id":"hfj***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088102179843724","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"003","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2019-10-22 15:59:49","total_amount":"100.00","trade_no":"2019102222001443721000074028","trade_status":"TRADE_CLOSED"}`, content)
	assert.Equal(t, "UpBhxKyuxGEMRethWKXj13Vj3rG7YJb+tC8YQdISwwerv7Hnr87O/hpSpEOMj30IrVfrUs4IMaaHHI+VaNM9WlmdqkJMXcqrVkh0VTzLykYWsfI08EbYMkbj0MyNFO27JjbALq7kknFzjj2p3lf6X+zo1nE6Wh678yFZ3esVUIkI7YxB0sQQvxkG/kMefwDhGwigtljXf6QEylzKGHDQy9hlCOUHFDSM5XFd1DJQuFqhttG/ujOkRTSJy6DaV/kjjxaA0iLhDoCqcGHBxTN9WQtYV3I+hIHhOLUtmLmQzHTypA7bOOgKGv6zL+3zObftEIPdfN/pii2D5k6vbtmSbg==", sign)

	body = `{"alipay_trade_query_response":{"code":"10000","msg":"Success","buyer_logon_id":"hfj***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088102179843724","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"003","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2019-10-22 15:59:49","total_amount":"100.00","trade_no":"2019102222001443721000074028","trade_status":"TRADE_CLOSED"},"sign":"UpBhxKyuxGEMRethWKXj13Vj3rG7YJb+tC8YQdISwwerv7Hnr87O/hpSpEOMj30IrVfrUs4IMaaHHI+VaNM9WlmdqkJMXcqrVkh0VTzLykYWsfI08EbYMkbj0MyNFO27JjbALq7kknFzjj2p3lf6X+zo1nE6Wh678yFZ3esVUIkI7YxB0sQQvxkG/kMefwDhGwigtljXf6QEylzKGHDQy9hlCOUHFDSM5XFd1DJQuFqhttG/ujOkRTSJy6DaV/kjjxaA0iLhDoCqcGHBxTN9WQtYV3I+hIHhOLUtmLmQzHTypA7bOOgKGv6zL+3zObftEIPdfN/pii2D5k6vbtmSbg=="}`
	content, sign = parseJSONSource(body, rootNodeName)
	assert.Equal(t, `{"code":"10000","msg":"Success","buyer_logon_id":"hfj***@sandbox.com","buyer_pay_amount":"0.00","buyer_user_id":"2088102179843724","buyer_user_type":"PRIVATE","invoice_amount":"0.00","out_trade_no":"003","point_amount":"0.00","receipt_amount":"0.00","send_pay_date":"2019-10-22 15:59:49","total_amount":"100.00","trade_no":"2019102222001443721000074028","trade_status":"TRADE_CLOSED"}`, content)
	assert.Equal(t, "UpBhxKyuxGEMRethWKXj13Vj3rG7YJb+tC8YQdISwwerv7Hnr87O/hpSpEOMj30IrVfrUs4IMaaHHI+VaNM9WlmdqkJMXcqrVkh0VTzLykYWsfI08EbYMkbj0MyNFO27JjbALq7kknFzjj2p3lf6X+zo1nE6Wh678yFZ3esVUIkI7YxB0sQQvxkG/kMefwDhGwigtljXf6QEylzKGHDQy9hlCOUHFDSM5XFd1DJQuFqhttG/ujOkRTSJy6DaV/kjjxaA0iLhDoCqcGHBxTN9WQtYV3I+hIHhOLUtmLmQzHTypA7bOOgKGv6zL+3zObftEIPdfN/pii2D5k6vbtmSbg==", sign)
}
