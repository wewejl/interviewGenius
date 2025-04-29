package service

import (
	"errors"
	"fmt"
	"interviewGenius/internal/dto"
	"interviewGenius/internal/model"
	"interviewGenius/internal/pkg/payment"
	"net/http"

	"github.com/google/uuid"
)

type PaymentService struct {
	alipayService *payment.AlipayService
}

func NewPaymentService(config *payment.AlipayConfig) (*PaymentService, error) {
	alipayService, err := payment.NewAlipayService(config)
	if err != nil {
		return nil, err
	}

	return &PaymentService{
		alipayService: alipayService,
	}, nil
}

// CreatePayment 创建支付订单
func (s *PaymentService) CreatePayment(userID string, cardID uint) (*dto.PaymentResponse, error) {
	// 获取会员卡信息
	card, err := model.GetMemberCardByID(cardID)
	if err != nil {
		return nil, errors.New("会员卡不存在")
	}

	// 将userID字符串转换为UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("用户ID格式无效")
	}

	// 创建订单
	order, err := model.CreateOrder(userUUID, cardID, card.Price)
	if err != nil {
		return nil, err
	}

	// 将价格从分转换为元，并格式化为字符串
	amountStr := fmt.Sprintf("%.2f", float64(card.Price)/100)

	// 创建支付宝支付
	payURL, err := s.alipayService.CreatePayment(
		order.ID.String(),
		amountStr,
		card.Name,
	)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		OrderID: order.ID.String(),
		PayURL:  payURL,
	}, nil
}

// HandlePaymentNotify 处理支付通知
func (s *PaymentService) HandlePaymentNotify(params map[string]string) error {
	// 创建一个http.Request用于验证通知
	req := &http.Request{}

	// 验证支付宝通知
	if err := s.alipayService.VerifyNotification(req); err != nil {
		return err
	}

	// 获取订单ID
	orderID, err := s.GetOrderID(params)
	if err != nil {
		return err
	}

	// 将orderID字符串转换为UUID
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errors.New("订单ID格式无效")
	}

	// 更新订单状态
	return model.PayOrder(orderUUID)
}

// HandlePaymentReturn 处理支付返回
func (s *PaymentService) HandlePaymentReturn(params map[string]string) error {
	// 验证签名（这里依赖于支付宝SDK的验证方法）
	req := &http.Request{}
	if err := s.alipayService.VerifyNotification(req); err != nil {
		return errors.New("签名验证失败")
	}

	// 获取订单ID
	orderID, err := s.GetOrderID(params)
	if err != nil {
		return errors.New("获取订单ID失败")
	}

	// 将orderID字符串转换为UUID
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errors.New("订单ID格式无效")
	}

	// 更新订单状态
	if err := model.PayOrder(orderUUID); err != nil {
		return errors.New("更新订单状态失败")
	}

	return nil
}

// GetOrderID 从支付参数中获取订单ID
func (s *PaymentService) GetOrderID(params map[string]string) (string, error) {
	orderID, ok := params["out_trade_no"]
	if !ok {
		return "", errors.New("订单ID不存在")
	}
	return orderID, nil
}

// VerifySign 验证签名
func (s *PaymentService) VerifySign(params map[string]string) error {
	// 创建一个http.Request用于验证签名
	req := &http.Request{}
	return s.alipayService.VerifyNotification(req)
}
