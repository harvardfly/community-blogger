package ginalipay

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

/**
 ** APP跟wap支付合集，内部会对入参做统一的检验跟格式化操作
 */
func (a *AliPay) Pay(ctx context.Context, args *PayRequest) (*PayResponse, error) {
	if args == nil {
		return nil, errors.New("empty args")
	}

	if args.TotalAmount <= 0 {
		return nil, errors.New("invalid total_amount")
	}

	if args.Subject == "" {
		return nil, errors.New("invalid subject")
	}

	if args.OutTradeNo == "" {
		return nil, errors.New("empty out_trade_no")
	}

	if args.ProductCode != "0" && args.ProductCode != "1" {
		return nil, errors.New("invalid product_code")
	}

	if args.PayType == Pay_TYPE_WAP && args.QuitURL == "" {
		return nil, errors.New("wap_pay must enter quit_url")
	}

	// 该笔订单允许的最晚付款时间
	tos := int64(args.TimeoutExpire.Sub(time.Now()).Seconds())
	tom := tos / 60
	if tos%60 >= 50 {
		tom++
	}
	if tom < 1 {
		tom = 1
	}
	timeoutExpress := fmt.Sprintf("%dm", tom)
	// 绝对超时时间
	timeExpire := args.TimeoutExpire.Format("2006-01-02 15:04")
	subject := TruncateString(args.Subject, 256)
	body := TruncateString(args.Body, 128)
	totalAmount := fmt.Sprintf("%.2f", float64(args.TotalAmount)/100.0)
	outTradeNo, err := MustLength(args.OutTradeNo, 64)
	if err != nil {
		return nil, err
	}

	passbackParams, err := MustLength(args.PassbackParams, 512)
	if err != nil {
		return nil, err
	}

	notifyURL, err := MustLength(args.NotifyURL, 256)
	if err != nil {
		return nil, err
	}

	quitURL, err := MustLength(args.QuitURL, 400)
	if err != nil {
		return nil, err
	}

	if args.PayType == Pay_TYPE_WAP {
		wapRes, err := a.TradeWapPay(ctx, &TradeWapPayRequest{
			TotalAmount:    totalAmount,
			Subject:        subject,
			OutTradeNo:     outTradeNo,
			ProductCode:    args.ProductCode,
			QuitURL:        quitURL,
			Body:           body,
			NotifyURL:      notifyURL,
			PassbackParams: passbackParams,
			TimeoutExpress: timeoutExpress,
			TimeExpire:     timeExpire,
		})
		if err != nil {
			return nil, err
		}

		return &PayResponse{
			WapPayTargetURL: wapRes.TargetURL,
		}, nil
	}

	appRes, err := a.TradeAppPay(ctx, &TradeAppPayRequest{
		TotalAmount:    totalAmount,
		Subject:        subject,
		OutTradeNo:     outTradeNo,
		ProductCode:    args.ProductCode,
		Body:           body,
		NotifyURL:      notifyURL,
		PassbackParams: passbackParams,
		TimeoutExpress: timeoutExpress,
		TimeExpire:     timeExpire,
	})
	if err != nil {
		return nil, err
	}

	return &PayResponse{
		AppPayBody: appRes.Body,
	}, nil
}

/**
 ** app支付接口2.0
 * 外部商户APP唤起快捷SDK创建订单并支付
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.app.pay
 */
func (a *AliPay) TradeAppPay(ctx context.Context, args *TradeAppPayRequest) (*TradeAppPayResponse, error) {
	p, err := a.URLValues(args)
	if err != nil {
		return nil, err
	}

	return &TradeAppPayResponse{
		Body: p.Encode(),
	}, nil
}

/**
 ** 手机网站支付接口2.0
 *
 * 外部商户创建订单并支付
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.wap.pay
 */
func (a *AliPay) TradeWapPay(ctx context.Context, args *TradeWapPayRequest) (*TradeWapPayResponse, error) {
	p, err := a.URLValues(args)
	if err != nil {
		return nil, err
	}

	return &TradeWapPayResponse{
		TargetURL: a.apiDomain + "?" + p.Encode(),
	}, nil
}

/**
 ** 统一收单交易关闭接口
 *
 * 用于交易创建后，用户在一定时间内未进行支付，可调用该接口直接将未付款的交易进行关闭。
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.close
 */
func (a *AliPay) TradeClose(ctx context.Context, args *TradeCloseRequest) (result *TradeCloseResponse, err error) {
	result = &TradeCloseResponse{}
	err = a.doRequestWithVerify(http.MethodPost, args, result)
	return
}

/**
 **
 ** 支付交易返回失败或支付系统超时，调用该接口撤销交易。如果此订单用户支付失败，支付宝系统会将此订单关闭；如果用户支付成功，支付宝系统会将此订单资金退还给用户。 注意：只有发生支付系统超时或者支付结果未知时可调用撤销，其他正常支付的单如需实现相同功能请调用申请退款API。提交支付交易后调用【查询订单API】，没有明确的支付结果再调用【撤销订单API】。
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.cancel
 */
func (a *AliPay) TradeCancel(ctx context.Context, args *TradeCancelRequest) (result *TradeCancelResponse, err error) {
	result = &TradeCancelResponse{}
	err = a.doRequestWithVerify(http.MethodPost, args, result)
	return
}

/**
 ** 交易查询接口
 *
 * 该接口提供所有支付宝支付订单的查询，商户可以通过该接口主动查询订单状态，完成下一步的业务逻辑
 * 需要调用查询接口的情况： 当商户后台、网络、服务器等出现异常，商户系统最终未接收到支付通知；
 * 调用支付接口后，返回系统错误或未知交易状态情况； 调用alipay.trade.pay，返回INPROCESS的状态； 调用alipay.trade.cancel之前，需确认支付状态；
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.query
 */
func (a *AliPay) TradeQuery(ctx context.Context, args *TradeQueryRequest) (result *TradeQueryResponse, err error) {
	result = &TradeQueryResponse{}
	err = a.doRequestWithVerify(http.MethodPost, args, result)
	return
}

/**
 ** 交易退款接口
 *
 * 当交易发生之后一段时间内，由于买家或者卖家的原因需要退款时，卖家可以通过退款接口将支付款退还给买家，支付宝将在收到退款请求并且验证成功之后，
 * 按照退款规则将支付款按原路退到买家帐号上。 交易超过约定时间（签约时设置的可退款时间）的订单无法进行退款 支付宝退款支持单笔交易分多次退款，
 * 多次退款需要提交原支付订单的商户订单号和设置不同的退款单号。一笔退款失败后重新提交，要采用原来的退款单号。总退款金额不能超过用户实际支付金额
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.refund
 */
func (a *AliPay) TradeRefund(ctx context.Context, args *TradeRefundRequest) (result *TradeRefundResponse, err error) {
	result = &TradeRefundResponse{}
	err = a.doRequestWithVerify(http.MethodPost, args, result)
	return
}

/**
 ** 查询账单下载地址接口
 *
 *  商户可使用该接口查询自已通过alipay.trade.refund或alipay.trade.refund.apply提交的退款请求是否执行成功。
 *  该接口的返回码10000，仅代表本次查询操作成功，不代表退款成功。如果该接口返回了查询数据，且refund_status为空或为REFUND_SUCCESS，则代表退款成功，
 *  如果没有查询到则代表未退款成功，可以调用退款接口进行重试。重试时请务必保证退款请求号一致。
 *
 * https://docs.open.alipay.com/api_1/alipay.trade.fastpay.refund.query
 */
func (a *AliPay) TradeRefundQuery(ctx context.Context, args *TradeRefundQueryRequest) (result *TradeRefundQueryResponse, err error) {
	result = &TradeRefundQueryResponse{}
	err = a.doRequestWithVerify(http.MethodPost, args, result)
	return
}

/**
 ** 查询账单下载地址接口
 *
 * 为方便商户快速查账，支持商户通过本接口获取商户离线账单下载地址
 *
 * https://docs.open.alipay.com/api_15/alipay.data.dataservice.bill.downloadurl.query
 */
func (a *AliPay) BillDownloadURLQuery(ctx context.Context, args *BillDownloadRequest) (result *BillDownloadResponse, err error) {
	result = &BillDownloadResponse{}
	err = a.doRequestWithVerify(http.MethodPost, args, result)
	return
}
