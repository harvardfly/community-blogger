package alipay

const (
	SandboxApiURL    = "https://openapi.alipaydev.com/gateway.do"
	ProductionApiURL = "https://openapi.alipay.com/gateway.do"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	FORMAT     = "JSON"
	CHARSET    = "utf-8"
	VERSION    = "1.0"
)

const (
	ResponseSuffix = "_response"
)

const (
	SignTypeRsa2 = "RSA2"
	SignTypeRsa  = "RSA"
)

type PayType int

const (
	PayTypeApp  PayType = iota // 手机支付
	PayTypeWap                 // 手机网站支付
	PayTypePage                // pc网站支付
)

// 交易状态说明(trade_status)
const (
	TradeStatusWaitBuyerPay  = "WAIT_BUYER_PAY" // 交易创建，等待买家付款
	TradeStatusTradeClosed   = "TRADE_CLOSED"   // 未付款交易超时关闭，或支付完成后全额退款
	TradeStatusTradeSuccess  = "TRADE_SUCCESS"  // 交易支付成功
	TradeStatusTradeFinished = "TRADE_FINISHED" // 交易结束，不可退款
)

// 结算操作类型(OperationType)
const (
	OperationTypeReplenish       = "replenish"
	OperationTypeReplenishRefund = "replenish_refund"
	OperationTypeTransfer        = "transfer"
	OperationTypeTransferRefund  = "transfer_refund"
)

// 买家用户类型
const (
	BuyerUserTypeCorporate = "CORPORATE" // 企业用户
	BuyerUserTypePrivate   = "PRIVATE"   // 个人用户
)

// 账单类型
const (
	BillTypeTrade        = "trade"        // 指商户基于支付宝交易收单的业务账单
	BillTypeSignCustomer = "signcustomer" // 是指基于商户支付宝余额收入及支出等资金变动的帐务账单
)

// 分账类型
const (
	RoyaltyTypeTransfer  = "transfer"  // 普通分账
	RoyaltyTypeReplenish = "replenish" // 补差
)

// 渠道所使用的资金类型
const (
	FundTypeDebiCard   = "DEBIT_CARD"  // 借记卡
	FundTypeCreditCard = "CREDIT_CARD" // 信用卡
	FundTypeMixedCard  = "MIXED_CARD"  // 借贷合一卡
)
