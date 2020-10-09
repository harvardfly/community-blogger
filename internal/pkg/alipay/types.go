package alipay

import (
	"time"
)

type CommonResponse struct {
	Code    string `json:"code"`     // 网关返回码
	Msg     string `json:"msg"`      // 网关返回码描述
	SubCode string `json:"sub_code"` // 可选：业务返回码，参见具体的API接口文档
	SubMsg  string `json:"sub_msg"`  // 可选：业务返回码描述，参见具体的API接口文档
}

type PayRequest struct {
	PayType        PayType   `json:"pay_type"`        // 支付类型
	TotalAmount    uint32    `json:"total_amount"`    // 支付总金额(单位分)
	Subject        string    `json:"subject"`         // 商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo     string    `json:"out_trade_no"`    // 订单号
	Body           string    `json:"body"`            // 商品描述
	NotifyURL      string    `json:"notify_url"`      // 回调地址
	TimeoutExpire  time.Time `json:"timeout_expire"`  // 该笔订单允许的最晚付款时间，逾期将关闭交易
	ProductCode    string    `json:"product_code"`    // 销售产品码，商家和支付宝签约的产品码
	PassbackParams string    `json:"passback_params"` // 公用回传参数
	GoodsType      string    `json:"goods_type"`      // 商品主类型 :0-虚拟类商品,1-实物类商品
	QuitURL        string    `json:"quit_url"`        // 用户付款中途退出返回商户网站的地址(如果是手机网站支付该参数必传)
}

type PayResponse struct {
	AppPayBody      string `json:"body"`       // 手机支付body体
	WapPayTargetURL string `json:"target_url"` // 手机网站支付跳转URL
}

/**
 * alipay.trade.app.pay
 */
type TradeAppPayRequest struct {
	TotalAmount string `json:"total_amount"` // 必选：订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	Subject     string `json:"subject"`      // 必选：商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo  string `json:"out_trade_no"` // 必选：商户网站唯一订单号

	AppAuthToken        string                   `json:"-"`                     // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	ReturnURL           string                   `json:"-"`                     // 可选：HTTP/HTTPS开头字符串
	NotifyURL           string                   `json:"-"`                     // 可选：支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	TimeoutExpress      string                   `json:"timeout_express"`       // 可选：该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	ProductCode         string                   `json:"product_code"`          // 可选：销售产品码，商家和支付宝签约的产品码
	Body                string                   `json:"body"`                  // 可选：对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body
	TimeExpire          string                   `json:"time_expire"`           // 可选：绝对超时时间，格式为yyyy-MM-dd HH:mm。
	GoodsType           string                   `json:"goods_type"`            // 可选：商品主类型 :0-虚拟类商品,1-实物类商品
	PromoParams         string                   `json:"promo_params"`          // 可选：优惠参数 （仅与支付宝协商后可用）
	PassbackParams      string                   `json:"passback_params"`       // 可选：公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在同步返回（包括跳转回商户网站）和异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝。
	ExtendParams        *TradeAppPayExtendParams `json:"extend_params"`         // 可选：业务扩展参数
	EnablePayChannels   string                   `json:"enable_pay_channels"`   // 可选：可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 ，与disable_pay_channels互斥
	StoreID             string                   `json:"store_id"`              // 可选：商户门店编号
	SpecifiedChannel    string                   `json:"specified_channel"`     // 可选：指定渠道，目前仅支持传入pcredit 若由于用户原因渠道不可用，用户可选择是否用其他渠道支付。 该参数不可与花呗分期参数同时传入
	DisablePayChannels  string                   `json:"disable_pay_channels"`  // 可选：禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 ，与enable_pay_channels互斥
	BusinessParams      string                   `json:"business_params"`       // 可选：商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	GoodsDetail         interface{}              `json:"goods_detail"`          // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
	ExtUserInfo         interface{}              `json:"ext_user_info"`         // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
	AgreementSignParams interface{}              `json:"agreement_sign_params"` // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
}

func (t *TradeAppPayRequest) Method() string {
	return "alipay.trade.app.pay"
}

func (t *TradeAppPayRequest) Params() map[string]string {
	return map[string]string{
		"return_url": t.ReturnURL,
		"notify_url": t.NotifyURL,
		//"app_auth_token": t.AppAuthToken,
	}
}

func (t *TradeAppPayRequest) BizContent() string {
	return marshal(t)
}

type TradeAppPayExtendParams struct {
	SysServiceProviderID string `json:"sys_service_provider_id"` // 可选：系统商编号 该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID
	HBFQNum              string `json:"hb_fq_num"`               // 可选：使用花呗分期要进行的分期数
	HBFQSellerPercent    string `json:"hb_fq_seller_percent"`    // 可选：使用花呗分期需要卖家承担的手续费比例的百分值，传入100代表100%
	IndustryRefluxInfo   string `json:"industry_reflux_info"`    // 可选：行业数据回流信息, 详见：地铁支付接口参数补充说明
	CardType             string `json:"card_type"`               // 可选：卡类型
}

type TradeAppPayResponse struct {
	Body string `json:"body"`
}

/**
 * alipay.trade.wap.pay
 */
type TradeWapPayRequest struct {
	TotalAmount string `json:"total_amount"` // 必选：订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	Subject     string `json:"subject"`      // 必选：商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo  string `json:"out_trade_no"` // 必选：商户网站唯一订单号
	ProductCode string `json:"product_code"` // 必选：销售产品码，商家和支付宝签约的产品码
	QuitURL     string `json:"quit_url"`     // 必选：用户付款中途退出返回商户网站的地址

	AppAuthToken        string                   `json:"-"`                     // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	ReturnURL           string                   `json:"-"`                     // 可选：HTTP/HTTPS开头字符串
	NotifyURL           string                   `json:"-"`                     // 可选：支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	TimeoutExpress      string                   `json:"timeout_express"`       // 可选：该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	Body                string                   `json:"body"`                  // 可选：对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body
	TimeExpire          string                   `json:"time_expire"`           // 可选：绝对超时时间，格式为yyyy-MM-dd HH:mm。
	GoodsType           string                   `json:"goods_type"`            // 可选：商品主类型 :0-虚拟类商品,1-实物类商品
	PromoParams         string                   `json:"promo_params"`          // 可选：优惠参数 （仅与支付宝协商后可用）
	PassbackParams      string                   `json:"passback_params"`       // 可选：公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在同步返回（包括跳转回商户网站）和异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝。
	ExtendParams        *TradeAppPayExtendParams `json:"extend_params"`         // 可选：业务扩展参数
	EnablePayChannels   string                   `json:"enable_pay_channels"`   // 可选：可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 ，与disable_pay_channels互斥
	StoreID             string                   `json:"store_id"`              // 可选：商户门店编号
	SpecifiedChannel    string                   `json:"specified_channel"`     // 可选：指定渠道，目前仅支持传入pcredit 若由于用户原因渠道不可用，用户可选择是否用其他渠道支付。 该参数不可与花呗分期参数同时传入
	DisablePayChannels  string                   `json:"disable_pay_channels"`  // 可选：禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 ，与enable_pay_channels互斥
	BusinessParams      string                   `json:"business_params"`       // 可选：商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	GoodsDetail         interface{}              `json:"goods_detail"`          // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
	ExtUserInfo         interface{}              `json:"ext_user_info"`         // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
	AgreementSignParams interface{}              `json:"agreement_sign_params"` // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
}

func (t *TradeWapPayRequest) Method() string {
	return "alipay.trade.wap.pay"
}

func (t *TradeWapPayRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = t.NotifyURL
	m["return_url"] = t.ReturnURL
	return m
}

func (t *TradeWapPayRequest) BizContent() string {
	return marshal(t)
}

type TradeWapPayResponse struct {
	TargetURL string `json:"target_url"`
}

/* alipay.trade.page.pay
 */
type TradePagePayRequest struct {
	TotalAmount string `json:"total_amount"` // 必选：订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	Subject     string `json:"subject"`      // 必选：商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo  string `json:"out_trade_no"` // 必选：商户网站唯一订单号
	ProductCode string `json:"product_code"` // 必选：销售产品码，商家和支付宝签约的产品码

	AppAuthToken        string                   `json:"-"`                     // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	ReturnURL           string                   `json:"-"`                     // 可选：HTTP/HTTPS开头字符串
	NotifyURL           string                   `json:"-"`                     // 可选：支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	TimeoutExpress      string                   `json:"timeout_express"`       // 可选：该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	Body                string                   `json:"body"`                  // 可选：对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body
	TimeExpire          string                   `json:"time_expire"`           // 可选：绝对超时时间，格式为yyyy-MM-dd HH:mm。
	GoodsType           string                   `json:"goods_type"`            // 可选：商品主类型 :0-虚拟类商品,1-实物类商品
	PromoParams         string                   `json:"promo_params"`          // 可选：优惠参数 （仅与支付宝协商后可用）
	PassbackParams      string                   `json:"passback_params"`       // 可选：公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在同步返回（包括跳转回商户网站）和异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝。
	ExtendParams        *TradeAppPayExtendParams `json:"extend_params"`         // 可选：业务扩展参数
	EnablePayChannels   string                   `json:"enable_pay_channels"`   // 可选：可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 ，与disable_pay_channels互斥
	StoreID             string                   `json:"store_id"`              // 可选：商户门店编号
	DisablePayChannels  string                   `json:"disable_pay_channels"`  // 可选：禁用渠道，用户不可用指定渠道支付 当有多个渠道时用“,”分隔 ，与enable_pay_channels互斥
	BusinessParams      string                   `json:"business_params"`       // 可选：商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	GoodsDetail         interface{}              `json:"goods_detail"`          // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
	ExtUserInfo         interface{}              `json:"ext_user_info"`         // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
	AgreementSignParams interface{}              `json:"agreement_sign_params"` // 可选，暂时用不到没有定义，需要时添加定义然后赋值即可
}

func (t *TradePagePayRequest) Method() string {
	return "alipay.trade.page.pay"
}

func (t *TradePagePayRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["notify_url"] = t.NotifyURL
	m["return_url"] = t.ReturnURL
	return m
}

func (t *TradePagePayRequest) BizContent() string {
	return marshal(t)
}

type TradePagePayResponse struct {
	TargetURL string `json:"target_url"`
}

/**
 * alipay.trade.close
 */
type TradeCloseRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 特殊可选：订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo    string `json:"trade_no"`     // 特殊可选：支付宝交易号，和商户订单号不能同时为空

	AppAuthToken string `json:"-"`           // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	OperatorID   string `json:"operator_id"` // 可选：卖家端自定义的的操作员 ID
}

func (t *TradeCloseRequest) Method() string {
	return "alipay.trade.close"
}

func (t *TradeCloseRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

func (t *TradeCloseRequest) BizContent() string {
	return marshal(t)
}

type TradeCloseResponse struct {
	TradeCloseResponse struct {
		CommonResponse
		TradeNo    string `json:"trade_no"`     // 支付宝交易号
		OutTradeNo string `json:"out_trade_no"` // 商家订单号
	} `json:"alipay_trade_close_response"`
	Sign string `json:"sign"` // 签名
}

/**
 * alipay.trade.cancel
 */
type TradeCancelRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 特殊可选：订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo    string `json:"trade_no"`     // 特殊可选：支付宝交易号，和商户订单号不能同时为空

	AppAuthToken string `json:"-"` // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
}

func (t *TradeCancelRequest) Method() string {
	return "alipay.trade.cancel"
}

func (t *TradeCancelRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

func (t *TradeCancelRequest) BizContent() string {
	return marshal(t)
}

type TradeCancelResponse struct {
	TradeCancelResponse struct {
		CommonResponse
		TradeNo            string    `json:"trade_no"`             // 支付宝交易号;当发生交易关闭或交易退款时返回；
		OutTradeNo         string    `json:"out_trade_no"`         // 商家订单号
		RetryFlag          string    `json:"retry_flag"`           // 是否需要重试
		GMTRefundPay       time.Time `json:"gmt_refund_pay"`       // 当撤销产生了退款时，返回退款时间； 默认不返回该信息，需与支付宝约定后配置返回；
		RefundSettlementID string    `json:"refund_settlement_id"` // 当撤销产生了退款时，返回的退款清算编号，用于清算对账使用； 只在银行间联交易场景下返回该信息；
	} `json:"alipay_trade_cancel_response"`
	Sign string `json:"sign"` // 签名
}

/**
 * alipay.trade.query
 */
type TradeQueryRequest struct {
	OutTradeNo string `json:"out_trade_no"` // 特殊可选：订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo    string `json:"trade_no"`     // 特殊可选：支付宝交易号，和商户订单号不能同时为空

	AppAuthToken string   `json:"-"`             // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	OrgPID       string   `json:"org_pid"`       // 可选：银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要查询的交易所属收单机构的pid;
	QueryOptions []string `json:"query_options"` // 可选：查询选项，商户通过上送该字段来定制查询返回信息
}

func (t *TradeQueryRequest) Method() string {
	return "alipay.trade.query"
}

func (t *TradeQueryRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

func (t *TradeQueryRequest) BizContent() string {
	return marshal(t)
}

type TradeQueryResponse struct {
	TradeQueryResponse struct {
		CommonResponse
		TradeNo             string           `json:"trade_no"`               // 支付宝交易号
		OutTradeNo          string           `json:"out_trade_no"`           // 商家订单号
		BuyerLogonID        string           `json:"buyer_logon_id"`         // 买家支付宝账号
		TradeStatus         string           `json:"trade_status"`           // 交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
		TotalAmount         string           `json:"total_amount"`           // 交易的订单金额，单位为元，两位小数。该参数的值为支付时传入的total_amount
		TransCurrency       string           `json:"trans_currency"`         // 标价币种，该参数的值为支付时传入的trans_currency，支持英镑：GBP、港币：HKD、美元：USD、新加坡元：SGD、日元：JPY、加拿大元：CAD、澳元：AUD、欧元：EUR、新西兰元：NZD、韩元：KRW、泰铢：THB、瑞士法郎：CHF、瑞典克朗：SEK、丹麦克朗：DKK、挪威克朗：NOK、马来西亚林吉特：MYR、印尼卢比：IDR、菲律宾比索：PHP、毛里求斯卢比：MUR、以色列新谢克尔：ILS、斯里兰卡卢比：LKR、俄罗斯卢布：RUB、阿联酋迪拉姆：AED、捷克克朗：CZK、南非兰特：ZAR、人民币：CNY、新台币：TWD。当trans_currency 和 settle_currency 不一致时，trans_currency支持人民币：CNY、新台币：TWD
		SettleCurrency      string           `json:"settle_currency"`        // 订单结算币种，对应支付接口传入的settle_currency，支持英镑：GBP、港币：HKD、美元：USD、新加坡元：SGD、日元：JPY、加拿大元：CAD、澳元：AUD、欧元：EUR、新西兰元：NZD、韩元：KRW、泰铢：THB、瑞士法郎：CHF、瑞典克朗：SEK、丹麦克朗：DKK、挪威克朗：NOK、马来西亚林吉特：MYR、印尼卢比：IDR、菲律宾比索：PHP、毛里求斯卢比：MUR、以色列新谢克尔：ILS、斯里兰卡卢比：LKR、俄罗斯卢布：RUB、阿联酋迪拉姆：AED、捷克克朗：CZK、南非兰特：ZAR
		SettleAmount        float32          `json:"settle_amount"`          // 结算币种订单金额
		PayCurrency         string           `json:"pay_currency"`           // 订单支付币种
		PayAmount           string           `json:"pay_amount"`             // 支付币种订单金额
		SettleTransRate     string           `json:"settle_trans_rate"`      // 结算币种兑换标价币种汇率
		TransPayRate        string           `json:"trans_pay_rate"`         // 标价币种兑换支付币种汇率
		BuyerPayAmount      string           `json:"buyer_pay_amount"`       // 买家实付金额，单位为元，两位小数。该金额代表该笔交易买家实际支付的金额，不包含商户折扣等金额
		PointAmount         string           `json:"point_amount"`           // 积分支付的金额，单位为元，两位小数。该金额代表该笔交易中用户使用积分支付的金额，比如集分宝或者支付宝实时优惠等
		InvoiceAmount       string           `json:"invoice_amount"`         // 交易中用户支付的可开具发票的金额，单位为元，两位小数。该金额代表该笔交易中可以给用户开具发票的金额
		SendPayDate         string           `json:"send_pay_date"`          // 本次交易打款给卖家的时间
		ReceiptAmount       string           `json:"receipt_amount"`         // 实收金额，单位为元，两位小数。该金额为本笔交易，商户账户能够实际收到的金额
		StoreID             string           `json:"store_id"`               // 商户门店编号
		TerminalID          string           `json:"terminal_id"`            // 商户机具终端编号
		FundBillList        []*TradeFundBill `json:"fund_bill_list"`         // 交易支付使用的资金渠道
		StoreName           string           `json:"store_name"`             // 请求交易支付中的商户店铺的名称
		BuyerUserID         string           `json:"buyer_user_id"`          // 买家在支付宝的用户id
		ChargeAmount        string           `json:"charge_amount"`          // 该笔交易针对收款方的收费金额； 默认不返回该信息，需与支付宝约定后配置返回；
		ChargeFlags         string           `json:"charge_flags"`           // 费率活动标识，当交易享受活动优惠费率时，返回该活动的标识； 默认不返回该信息，需与支付宝约定后配置返回； 可能的返回值列表： 蓝海活动标识：bluesea_1
		SettlementID        string           `json:"settlement_id"`          // 支付清算编号，用于清算对账使用； 只在银行间联交易场景下返回该信息；
		TradeSettleInfo     *TradeSettleInfo `json:"trade_settle_info"`      // 返回的交易结算信息，包含分账、补差等信息
		AuthTradePayMode    string           `json:"auth_trade_pay_mode"`    // 预授权支付模式，该参数仅在信用预授权支付场景下返回。信用预授权支付：CREDIT_PREAUTH_PAY
		BuyerUserType       string           `json:"buyer_user_type"`        // 买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
		MdiscountAmount     string           `json:"mdiscount_amount"`       // 商家优惠金额
		DiscountAmount      string           `json:"discount_amount"`        // 平台优惠金额
		BuyerUserName       string           `json:"buyer_user_name"`        // 买家名称； 买家为个人用户时为买家姓名，买家为企业用户时为企业名称； 默认不返回该信息，需与支付宝约定后配置返回
		Subject             string           `json:"subject"`                // 订单标题； 只在间连场景下返回；
		Body                string           `json:"body"`                   // 订单描述; 只在间连场景下返回；
		AlipaySubMerchantID string           `json:"alipay_sub_merchant_id"` // 间连商户在支付宝端的商户编号； 只在间连场景下返回；
		ExtInfos            string           `json:"ext_infos"`              // 交易额外信息，特殊场景下与支付宝约定返回。 json格式。
	} `json:"alipay_trade_query_response"` // response
	Sign string `json:"sign"`              // 签名
}

type TradeSettleInfo struct {
	TradeSettleDetailList []*TradeSettleDetail `json:"trade_settle_detail_list"` //交易结算明细信息
}

type TradeSettleDetail struct {
	OperationType     string    `json:"operation_type"`      // 结算操作类型。包含replenish、replenish_refund、transfer、transfer_refund等类型
	OperationSerialNo string    `json:"operation_serial_no"` // 商户操作序列号。商户发起请求的外部请求号。
	OperationDT       time.Time `json:"operation_dt"`        // 操作日期
	TransOut          string    `json:"trans_out"`           // 转出账号
	TransIn           string    `json:"trans_in"`            // 转入账号
	Amount            float32   `json:"amount"`              // 实际操作金额，单位为元，两位小数。该参数的值为分账或补差或结算时传入
}

type TradeFundBill struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道
	BankCode    string `json:"bank_code"`    // 银行卡支付时的银行代码
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
	FundType    string `json:"fund_type"`    // 渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡)
}

/**
 * alipay.trade.refund
 */
type TradeRefundRequest struct {
	RefundAmount float32 `json:"refund_amount"` // 必选：需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	OutTradeNo   string  `json:"out_trade_no"`  // 特殊可选：订单支付时传入的商户订单号,不能和 trade_no同时为空。
	TradeNo      string  `json:"trade_no"`      // 特殊可选：支付宝交易号，和商户订单号不能同时为空

	AppAuthToken   string         `json:"-"`               // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	RefundCurrency string         `json:"refund_currency"` // 可选：订单退款币种信息
	RefundReason   string         `json:"refund_reason"`   // 可选：退款的原因说明
	OutRequestNo   string         `json:"out_request_no"`  // 可选：标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorID     string         `json:"operator_id"`     // 可选：商户的操作员编号
	StoreID        string         `json:"store_id"`        // 可选：商户的门店编号
	TerminalID     string         `json:"terminal_id"`     // 可选：商户的终端编号
	GoodsDetail    []*GoodsDetail `json:"goods_detail"`    // 可选： 退款包含的商品列表信息，Json格式。 其它说明详见：“商品明细说明”
	OrgPID         string         `json:"org_pid"`         // 可选：银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要退款的交易所属收单机构的pid;
}

type GoodsDetail struct {
	GoodsName string  `json:"goods_name"` // 必选：商品名称
	Quantity  int     `json:"quantity"`   // 必选：商品数量
	Price     float32 `json:"price"`      // 必选：商品单价，单位为元

	GoodsID        string `json:"goods_id"`         // 可选：商品的编号
	AliPayGoodsID  string `json:"ali_pay_goods_id"` // 可选：支付宝定义的统一商品编号
	GoodsCategory  string `json:"goods_category"`   // 可选：商品类目
	CategoriesTree string `json:"categories_tree"`  // 可选：商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割
	Body           string `json:"body"`             // 可选：商品描述信息
	ShowURL        string `json:"show_url"`         // 可选：商品的展示地址
}

func (t *TradeRefundRequest) Method() string {
	return "alipay.trade.refund"
}

func (t *TradeRefundRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

func (t *TradeRefundRequest) BizContent() string {
	return marshal(t)
}

type TradeRefundResponse struct {
	TradeRefundResponse struct {
		CommonResponse
		TradeNo                      string               `json:"trade_no"`                        // 支付宝交易号
		OutTradeNo                   string               `json:"out_trade_no"`                    // 商家订单号
		BuyerLogonID                 string               `json:"buyer_logon_id"`                  // 买家支付宝账号
		FundChange                   string               `json:"fund_change"`                     // 本次退款是否发生了资金变化
		RefundFee                    string               `json:"refund_fee"`                      // 退款总金额
		RefundCurrency               string               `json:"refund_currency"`                 // 退款币种信息
		GMTRefundPay                 string               `json:"gmt_refund_pay"`                  // 退款支付时间
		RefundDetailItemList         []*TradeFundBill     `json:"refund_detail_item_list"`         // 退款使用的资金渠道
		StoreName                    string               `json:"store_name"`                      // 交易在支付时候的门店名称
		BuyerUserID                  string               `json:"buyer_user_id"`                   // 买家在支付宝的用户id
		RefundPresetPaytoolList      []*PresetPayToolInfo `json:"refund_preset_paytool_list"`      // 退回的前置资产列表
		RefundSettlementID           string               `json:"refund_settlement_id"`            // 退款清算编号，用于清算对账使用； 只在银行间联交易场景下返回该信息；
		PresentRefundBuyerAmount     string               `json:"present_refund_buyer_amount"`     // 本次退款金额中买家退款金额
		PresentRefundDiscountAmount  string               `json:"present_refund_discount_amount"`  // 本次退款金额中平台优惠退款金额
		PresentRefundMdiscountAmount string               `json:"present_refund_mdiscount_amount"` // 本次退款金额中商家优惠退款金额
	} `json:"alipay_trade_refund_response"`
	Sign string `json:"sign"` // 签名
}

type PresetPayToolInfo struct {
	Amount         []float32 `json:"amount"`           // 前置资产金额
	AssertTypeCode string    `json:"assert_type_code"` // 前置资产类型编码，和收单支付传入的preset_pay_tool里面的类型编码保持一致。
}

/**
 * alipay.trade.fastpay.refund.query
 */
type TradeRefundQueryRequest struct {
	TradeNo      string `json:"trade_no"`       // 特殊可选：支付宝交易号，和商户订单号不能同时为空
	OutTradeNo   string `json:"out_trade_no"`   // 特殊可选：订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no
	OutRequestNo string `json:"out_request_no"` // 必选：请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号

	AppAuthToken string `json:"-"`       // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
	OrgPID       string `json:"org_pid"` // 可选：银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要查询的交易所属收单机构的pid;
}

func (t *TradeRefundQueryRequest) Method() string {
	return "alipay.trade.fastpay.refund.query"
}

func (t *TradeRefundQueryRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

func (t *TradeRefundQueryRequest) BizContent() string {
	return marshal(t)
}

type TradeRefundQueryResponse struct {
	RefundQueryResponse struct {
		CommonResponse
		TradeNo                      string               `json:"trade_no"`                        // 支付宝交易号
		OutTradeNo                   string               `json:"out_trade_no"`                    // 商家订单号
		OutRequestNo                 string               `json:"out_request_no"`                  // 本笔退款对应的退款请求号
		RefundReason                 string               `json:"refund_reason"`                   // 发起退款时，传入的退款原因
		TotalAmount                  string               `json:"total_amount"`                    // 该笔退款所对应的交易的订单金额
		RefundAmount                 string               `json:"refund_amount"`                   // 本次退款请求，对应的退款金额
		RefundRoyaltys               *RefundRoyaltyResult `json:"refund_royaltys"`                 // 退分账明细信息
		GMTRefundPay                 string               `json:"gmt_refund_pay"`                  // 退款时间； 默认不返回该信息，需与支付宝约定后配置返回；
		RefundDetailItemList         []*TradeFundBill     `json:"refund_detail_item_list"`         // 本次退款使用的资金渠道； 默认不返回该信息，需与支付宝约定后配置返回；
		SendBackFee                  string               `json:"send_back_fee"`                   // 本次商户实际退回金额； 默认不返回该信息，需与支付宝约定后配置返回；
		RefundSettlementID           string               `json:"refund_settlement_id"`            // 退款清算编号，用于清算对账使用； 只在银行间联交易场景下返回该信息；
		PresentRefundBuyerAmount     string               `json:"present_refund_buyer_amount"`     // 本次退款金额中买家退款金额
		PresentRefundDiscountAmount  string               `json:"present_refund_discount_amount"`  // 本次退款金额中平台优惠退款金额
		PresentRefundMdiscountAmount string               `json:"present_refund_mdiscount_amount"` // 本次退款金额中商家优惠退款金额
	} `json:"alipay_trade_fastpay_refund_query_response"`
	Sign string `json:"sign"` // 签名
}

type RefundRoyaltyResult struct {
	RefundAmount  string `json:"refund_amount"`   // 退分账金额
	RoyaltyType   string `json:"royalty_type"`    // 分账类型:分账类型. 普通分账为：transfer; 补差为：replenish; 为空默认为分账transfer;
	ResultCode    string `json:"result_code"`     // 退分账结果码
	TransOut      string `json:"trans_out"`       // 转出人支付宝账号对应用户ID
	TransOutEmail string `json:"trans_out_email"` // 转出人支付宝账号
	TransIn       string `json:"trans_in"`        // 转入人支付宝账号对应用户ID
	TransInEmail  string `json:"trans_in_email"`  // 转入人支付宝账号
}

/*
 * alipay.data.dataservice.bill.downloadurl.query
 */
type BillDownloadRequest struct {
	BillType string `json:"bill_type"` // 必选：账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型：trade、signcustomer；trade指商户基于支付宝交易收单的业务账单；signcustomer是指基于商户支付宝余额收入及支出等资金变动的帐务账单。
	BillDate string `json:"bill_date"` // 必选：账单时间：日账单格式为yyyy-MM-dd，最早可下载2016年1月1日开始的日账单；月账单格式为yyyy-MM，最早可下载2016年1月开始的月账单

	AppAuthToken string `json:"-"` // 可选：应用授权（https://docs.open.alipay.com/20160728150111277227/intro）
}

func (t *BillDownloadRequest) Method() string {
	return "alipay.data.dataservice.bill.downloadurl.query"
}

func (t *BillDownloadRequest) Params() map[string]string {
	var m = make(map[string]string)
	m["app_auth_token"] = t.AppAuthToken
	return m
}

func (t *BillDownloadRequest) BizContent() string {
	return marshal(t)
}

type BillDownloadResponse struct {
	AlipayDataDataserviceBillDownloadurlQueryResponse struct {
		Code            string `json:"code"`
		Msg             string `json:"msg"`
		BillDownloadURL string `json:"bill_download_url"`
	} `json:"alipay_data_dataservice_bill_downloadurl_query_response"`
	Sign string `json:"sign"` // 签名
}
