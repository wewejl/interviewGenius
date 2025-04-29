package dto

// PaymentRequest 支付请求
type PaymentRequest struct {
	CardID uint `json:"card_id" binding:"required"`
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	OrderID string `json:"order_id"`
	PayURL  string `json:"pay_url"`
}

// PaymentNotifyRequest 支付通知请求
type PaymentNotifyRequest struct {
	NotifyTime     string `form:"notify_time"`
	NotifyType     string `form:"notify_type"`
	NotifyID       string `form:"notify_id"`
	SignType       string `form:"sign_type"`
	Sign           string `form:"sign"`
	TradeNo        string `form:"trade_no"`
	OutTradeNo     string `form:"out_trade_no"`
	TradeStatus    string `form:"trade_status"`
	TotalAmount    string `form:"total_amount"`
	ReceiptAmount  string `form:"receipt_amount"`
	BuyerPayAmount string `form:"buyer_pay_amount"`
	SellerID       string `form:"seller_id"`
	AppID          string `form:"app_id"`
}
