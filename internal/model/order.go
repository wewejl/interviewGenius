package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"   // 已创建
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusFailed    OrderStatus = "failed"    // 支付失败
)

// Order 订单模型
type Order struct {
	ID           uuid.UUID   `json:"id" gorm:"type:char(36);primaryKey"`
	UserID       uuid.UUID   `json:"user_id" gorm:"type:char(36);not null;index"`
	CardID       uint        `json:"card_id" gorm:"not null"`
	Amount       int         `json:"amount" gorm:"not null"`               // 订单金额（单位：分）
	Status       OrderStatus `json:"status" gorm:"size:20;not null;index"` // 订单状态
	PurchaseTime *time.Time  `json:"purchase_time"`                        // 购买时间
	PaymentTime  *time.Time  `json:"payment_time"`                         // 支付时间
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Card         *MemberCard `json:"card,omitempty" gorm:"foreignKey:CardID"`
}

// BeforeCreate 创建前生成UUID
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	if o.Status == "" {
		o.Status = OrderStatusCreated
	}

	now := time.Now()
	o.PurchaseTime = &now

	return nil
}

// CreateOrder 创建订单
func CreateOrder(userID uuid.UUID, cardID uint, amount int) (*Order, error) {
	order := &Order{
		UserID: userID,
		CardID: cardID,
		Amount: amount,
		Status: OrderStatusCreated,
	}

	if err := DB.Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByID 根据ID获取订单
func GetOrderByID(id uuid.UUID) (*Order, error) {
	var order Order
	if err := DB.Preload("Card").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetUserOrders 获取用户的订单列表
func GetUserOrders(userID uuid.UUID) ([]*Order, error) {
	var orders []*Order
	if err := DB.Preload("Card").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// PayOrder 支付订单
func PayOrder(id uuid.UUID) error {
	now := time.Now()
	return DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		if err := tx.Model(&Order{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       OrderStatusPaid,
			"payment_time": now,
		}).Error; err != nil {
			return err
		}

		// 获取订单信息
		var order Order
		if err := tx.Preload("Card").First(&order, "id = ?", id).Error; err != nil {
			return err
		}

		// 获取用户信息
		var user User
		if err := tx.First(&user, "id = ?", order.UserID).Error; err != nil {
			return err
		}

		// 计算新的会员到期时间
		var expiryTime time.Time
		if user.MemberExpiry == nil || user.MemberExpiry.Before(now) {
			// 如果用户不是会员或会员已过期，从当前时间开始计算
			expiryTime = now.AddDate(0, 0, order.Card.DurationDays)
		} else {
			// 如果用户是会员且未过期，从当前到期时间开始叠加
			expiryTime = user.MemberExpiry.AddDate(0, 0, order.Card.DurationDays)
		}

		// 更新用户的会员到期时间
		if err := tx.Model(&User{}).Where("id = ?", order.UserID).Updates(map[string]interface{}{
			"member_expiry": expiryTime,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// CancelOrder 取消订单
func CancelOrder(id uuid.UUID) error {
	return DB.Model(&Order{}).Where("id = ?", id).Update("status", OrderStatusCancelled).Error
}
