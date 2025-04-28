package payment

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"net/http"
)

type AlipayConfig struct {
	AppID        string
	PrivateKey   string
	AliPublicKey string
	NotifyURL    string
	ReturnURL    string
	IsProduction bool
}

type AlipayService struct {
	client *alipay.Client
	config *AlipayConfig
}

func NewAlipayService(config *AlipayConfig) (*AlipayService, error) {
	client, err := alipay.New(config.AppID, config.PrivateKey, config.IsProduction)
	if err != nil {
		return nil, fmt.Errorf("创建支付宝客户端失败: %v", err)
	}

	err = client.LoadAliPayPublicKey(config.AliPublicKey)
	if err != nil {
		return nil, fmt.Errorf("加载支付宝公钥失败: %v", err)
	}

	return &AlipayService{
		client: client,
		config: config,
	}, nil
}

func (s *AlipayService) CreatePayment(orderID string, amount string, subject string) (string, error) {
	var p = alipay.TradePagePay{}
	p.NotifyURL = s.config.NotifyURL
	p.ReturnURL = s.config.ReturnURL
	p.Subject = subject
	p.OutTradeNo = orderID
	p.TotalAmount = amount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := s.client.TradePagePay(p)
	if err != nil {
		return "", fmt.Errorf("创建支付订单失败: %v", err)
	}

	return url.String(), nil
}

func (s *AlipayService) VerifyNotification(req *http.Request) error {
	notification, err := s.client.GetTradeNotification(req)
	if err != nil {
		return fmt.Errorf("验证通知失败: %v", err)
	}

	if notification.TradeStatus != "TRADE_SUCCESS" {
		return fmt.Errorf("交易未成功: %s", notification.TradeStatus)
	}

	return nil
}
