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
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4efPzSQ/KoBDVOpWJh8MUOzEoNVTSTXyvg2VPtT4xOPo9v/fvT22QLdKMOktuAx3HWroaDLYPqg641mA6ZrghR+EO/7eWxgH0ZfYoo8/ryXmp/k2sEv2l2bd91q++RdqBMDHp3tjHHVz94Dy6KRVkrNZbZMAImkh7m6EbhLlYWOGVQTcFq03PNC9eMolIPvhIk9bphlJUiXLwIjwBshTw7eJ0v/zNtDR42Owf6MbxBkiWBlWkXDlvgoow9i7tsv1UHaB0cO2kvJrKGh0xHTdGQ3GmJXqG+iLcnsU1Zvw75q8qTBQ84xpYw/2VQQw4Lo9CaZppFPOo9T1kbO1sxYKfwIDAQAB
-----END PUBLIC KEY-----`)

	privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAzKYmwgwffvLfgNh64EL8KHG0vb92qOe3fZohbfSlCjyF0JhH
i06iY3BpXcO3kFSUJU0I8K6tSLFm7t0eCKg+m4wwW7FDUKRLIbV4aB7MFeyESTVo
HHo9/g+5GEoDCOKbXycFGxez8SoUe5cbfv822OSgYCVBS1n+3WmtDS2z/qQwdmEq
pODB+3ADOSzzT7v2wjctJXK7VY5ojlzOMHun/ubNac9W2A7ZIM6/71GA5EE3KmMQ
BQfEvwOS45UUvrNm4n6UUz52XHkGCtV1E7PoZOyYXTyVkkLB/PBQcJ66nmrTrwBC
SS6Qy5JucuR71MyXwcmd49m9gLHMr5pfsGAKIQIDAQABAoIBABBjPbdMQTlpOXyX
2T8cHhUfBdor+tSLuaXVMdgcPxsSvaR2jUQah+ZumgFoKsrj+vkBsjWTx0yK6DM4
ga8vydOq5jRxfXJ3tYx9e6ba0Hzw86bpH+0n3M77c1b+lj+f5dE4zq/ctgwZ0ooT
D6Cbz0iPBkissw3VpxMT40eq4T8CwjL2w3DS1N9RWvOG6DVypRnwEO0BTTLD5jHd
kZlrxQ3ICHmXpcXVvH37al+P24b2uAl9PjaW0FImxM3dButrr3r2oPAHpqpTyyrn
zHQcomJ6e/DRNUh3N/81G7IpSl4CCR32tvCw65LICQE5dCf9X32nF5a4v8q0Oi6J
A1XfEgECgYEA6wyzFMfPZV8vaJX98MB5t/gbDjhYmr87pSwn6ObB+kq2CrW/tDdN
Cutxp2im3996F7wLGp3m+uuc/Y0cVhZhxBwWX7GwD8HNIfsnSCHcadsdFwE1nSCV
nYhfT5PFSXXxLI7ouSQTSU5UbBPx2IHhRNqu9Vx27qYrBc9+AiEKQNcCgYEA3uPH
h3HVxE919O+UrnHZ7QFFvLpC/tPtpKnT8jcfN2JhoEhZCDbUbHdRwurZsqMbdtXk
EhgNuDe84NGzDEtQUtdHoAN6j9+f6l5dbSTpIOxOBX26L2MEiNMli29WfUo2Lu4m
IDLFm5nxuh1EU6TLvUqvec5k1O3XtZHn+0XaFccCgYAbhPQelbpBeyB41T7TBiX6
FRFN2+j7zTH1h7LwgLvrSv3/SQI13lel1KUM3aLUCT0pNDn0ltIpRJav1OqhZaNy
q3svHwWnAqC6vsX9mwFMX3wLanfVerAproLCNWHe1PE0r1KuAnxDk+kscjVZjPNL
9XKQhY/jJw9Ycc+l/ipOJwKBgQDayqV2Y1v/lDCp+vPsOX4+lF0sYXqaQtaoKL0Q
quNNrpk+iUY8NfZXctkjiP2hyVKQWG3FBS+SgcQ6vB4SF2wFpaV9gWyyBkYn/fGf
zfe6hbwgz8YP9hbhaMMWGHjCDDMb5lIukShBEeCjXU9Q/Bey/Lk7zEpWahKw/UTG
906YyQKBgB2rxUUgYtkpiKXXYQH5pp9jI1uEWcU9iQEt6xJSIciJylB0JBdN99Oq
uJDIJ2wiuMrVKBQafuxW2k6P18WF+NCXo3eYzRnmrvcipCu97kDDdBf6TmTOEKm8
6lYnLNiNiQXurNnxVO9fOMhLFPh0oAZgfpsYk5PlclqPufh6i7Xi
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
		OutTradeNo: "005",
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
