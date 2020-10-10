package alipay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gotest.tools/assert"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGetTradeNotificationByBody(t *testing.T) {
	body := "app_id=2016082801813467&biz_content=%7B%22total_amount%22%3A%221%22%2C%22subject%22%3A%22%E6%B5%8B%E8%AF%95%E5%95%86%E5%93%81%22%2C%22out_trade_no%22%3A%2200000000000000001%22%2C%22product_code%22%3A%22QUICK_WAP_WAY%22%2C%22quit_url%22%3A%22http%3A%2F%2Fwww.baidu.com%22%2C%22return_url%22%3A%22%22%2C%22notify_url%22%3A%22%22%2C%22app_auth_token%22%3A%22%22%2C%22timeout_express%22%3A%22%22%2C%22body%22%3A%22%22%2C%22time_expire%22%3A%22%22%2C%22goods_type%22%3A%22%22%2C%22promo_params%22%3A%22%22%2C%22passback_params%22%3A%22%22%2C%22extend_params%22%3Anull%2C%22enable_pay_channels%22%3A%22%22%2C%22store_id%22%3A%22%22%2C%22specified_channel%22%3A%22%22%2C%22disable_pay_channels%22%3A%22%22%2C%22business_params%22%3A%22%22%2C%22goods_detail%22%3Anull%2C%22ext_user_info%22%3Anull%2C%22agreement_sign_params%22%3Anull%7D&charset=utf-8&format=JSON&method=alipay.trade.wap.pay&notify_url=&return_url=&sign=RPn2Eqiauq5kQeaQvrobuQy7LAv5XWDK5hQLmAyQdp6SQ0urz5KWId1VoM9T9kcg0VA1PD8t9eD28Cz0f3qeXFch3Voc86O7mf6GhXVYyFFCD%2BQQBd2%2Bol3JVxcW9QkItwOvzWFCt1FZwXJzCxfGfL8bDbt85Y6ReS9Fvt%2BHX5aSEHra4O8T3pUbxxCZWLDMba4wjzoKhL9UHz%2BFZDUFTX7qeyuXeWoSv9wfQNb1q1iIgXJ88XPpiTCpONkDgNWaY0BoITBM7O0guL8Gb14cdufbG0%2Fdc27DYjtwrzSsglqMauZXQ5CHoMpum879D3CO%2FPROug4CR8RH7BfRaEyAMw%3D%3D&sign_type=RSA2&timestamp=2019-10-21+10%3A17%3A19&version=1.0&notify_id=1111"
	noti, err := client.GetTradeNotificationByBody(body)
	_ = err
	assert.Equal(t, nil, err)
	t.Logf("%#v", *noti)
}

func TestGetTradeNotification(t *testing.T) {
	cli := NewAliPay(appID, publicKey, privateKey, false)
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(map[string]interface{}{
			"code": 0,
		})
		writer.Write(buf.Bytes())
	})

	mux.HandleFunc("/notify/wap/pay", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.Form.Get("order_id")
		if orderID == "" {
			orderID = fmt.Sprintf("%v", rand.Int63n(1000000000))
		}

		paymentAmount := r.Form.Get("payment_amount")
		intAmount, err := strconv.ParseInt(paymentAmount, 10, 64)
		if err != nil {
			r.Header.Set("Content-Type", "application/json;charset=utf-8")
			buf := bytes.NewBuffer(nil)
			json.NewEncoder(buf).Encode(map[string]interface{}{
				"code": 400,
				"msg":  "invalid payment_amount",
			})
			w.Write(buf.Bytes())
			return
		}

		ret, err := cli.TradeWapPay(context.Background(), &TradeWapPayRequest{
			TotalAmount:    fmt.Sprintf("%.2f", float64(intAmount)/100.0),
			Subject:        "测试商品",
			OutTradeNo:     orderID,
			ProductCode:    "QUICK_WAP_WAY",
			QuitURL:        "http://www.tsgbqx.cn",
			NotifyURL:      "",
			TimeoutExpress: "30m",
			GoodsType:      "1",
			TimeExpire:     time.Now().Add(time.Minute * 30).Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			r.Header.Set("Content-Type", "application/json;charset=utf-8")
			buf := bytes.NewBuffer(nil)
			json.NewEncoder(buf).Encode(map[string]interface{}{
				"code": 500,
				"msg":  err.Error(),
			})
			w.Write(buf.Bytes())
			return
		}
		http.Redirect(w, r, ret.TargetURL, http.StatusPermanentRedirect)
		return
	})

	mux.HandleFunc("/notify/alipay/pay.action", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.ParseForm()

		fmt.Println("body:", string(body))
		noti, err := cli.GetTradeNotificationByBody(string(body))
		if err != nil {
			fmt.Println("错误：", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Printf("notification:%#v", *noti)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Success"))
	})
	http.ListenAndServe(":4096", mux)
}
