package ginalipay

const (
	ALI_PAY_SANDBOX_API_URL    = "https://openapi.alipaydev.com/gateway.do"
	ALI_PAY_PRODUCTION_API_URL = "https://openapi.alipay.com/gateway.do"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
	FORMAT      = "JSON"
	CHARSET     = "utf-8"
	VERSION     = "1.0"
)

const (
	RESPONSE_SUFFIX = "_response"
)

const (
	SIGN_TYPE_RSA2 = "RSA2"
	SIGN_TYPE_RSA  = "RSA"
)

type PayType int

const (
	PAY_TYPE_APP PayType = iota // 手机支付
	Pay_TYPE_WAP                // 手机网站支付
)

// 交易状态说明(trade_status)
const (
	TRADE_STATUS_WAIT_BUYER_PAY = "WAIT_BUYER_PAY" // 交易创建，等待买家付款
	TRADE_STATUS_TRADE_CLOSED   = "TRADE_CLOSED"   // 未付款交易超时关闭，或支付完成后全额退款
	TRADE_STATUS_TRADE_SUCCESS  = "TRADE_SUCCESS"  // 交易支付成功
	TRADE_STATUS_TRADE_FINISHED = "TRADE_FINISHED" // 交易结束，不可退款
)

// 结算操作类型(operation_type)
const (
	OPERATION_TYPE_REPLENISH        = "replenish"
	OPERATION_TYPE_REPLENISH_REFUND = "replenish_refund"
	OPERATION_TYPE_TRANSFER         = "transfer"
	OPERATION_TYPE_TRANSFER_REFUND  = "transfer_refund"
)

// 买家用户类型
const (
	BUYER_USER_TYPE_CORPORATE = "CORPORATE" // 企业用户
	BUYER_USER_TYPE_PRIVATE   = "PRIVATE"   // 个人用户
)

// 账单类型
const (
	BILL_TYPE_TRADE        = "trade"        // 指商户基于支付宝交易收单的业务账单
	BILL_TYPE_SIGNCUSTOMER = "signcustomer" // 是指基于商户支付宝余额收入及支出等资金变动的帐务账单
)

// 分账类型
const (
	ROYALTY_TYPE_TRANSFER  = "transfer"  // 普通分账
	ROYALTY_TYPE_REPLENISH = "replenish" // 补差
)

// 渠道所使用的资金类型
const (
	FUND_TYPE_DEBI_CARD   = "DEBIT_CARD"  // 借记卡
	FUND_TYPE_CREDIT_CARD = "CREDIT_CARD" // 信用卡
	FUND_YPE_MIXED_CARD   = "MIXED_CARD"  // 借贷合一卡
)
