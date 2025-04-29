package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"interviewGenius/internal/model"
	"net/http"
)

// GetMemberInfo 获取用户会员信息
func GetMemberInfo(c *gin.Context) {
	// 获取当前用户ID
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 检查用户是否是会员
	isMember, expiryTime, err := model.IsMember(userID)
	if err != nil {
		zap.L().Error("获取会员信息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取会员信息失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取会员信息成功",
		"data": gin.H{
			"is_member":     isMember,
			"expiry_time":   expiryTime,
			"service_left":  isMember, // 会员可无限次使用
			"service_count": 1,        // 非会员每天1次
		},
	})
}

// GetMemberCards 获取会员卡列表
func GetMemberCards(c *gin.Context) {
	cards, err := model.GetAllMemberCards()
	if err != nil {
		zap.L().Error("获取会员卡列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取会员卡列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取会员卡列表成功",
		"data": gin.H{
			"cards": cards,
		},
	})
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CardID uint `json:"card_id" binding:"required"`
}

// CreateOrder 创建会员卡订单
func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取当前用户ID
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 获取会员卡信息
	card, err := model.GetMemberCardByID(req.CardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "会员卡不存在",
			"data": nil,
		})
		return
	}

	// 创建订单
	order, err := model.CreateOrder(userID, card.ID, card.Price)
	if err != nil {
		zap.L().Error("创建订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "创建订单失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "创建订单成功",
		"data": gin.H{
			"order_id": order.ID,
			"amount":   order.Amount,
		},
	})
}

// PayOrderRequest 支付订单请求
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

// PayOrder 支付会员卡订单
func PayOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "缺少订单ID",
			"data": nil,
		})
		return
	}

	var req PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取当前用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的订单ID",
			"data": nil,
		})
		return
	}

	// 这里应该有支付逻辑，为了简化示例，直接标记为已支付
	if err := model.PayOrder(orderUUID); err != nil {
		zap.L().Error("支付订单失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "支付订单失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "支付订单成功",
		"data": nil,
	})
}

// GetOrders 获取用户订单列表
func GetOrders(c *gin.Context) {
	// 获取当前用户ID
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	orders, err := model.GetUserOrders(userID)
	if err != nil {
		zap.L().Error("获取订单列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取订单列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取订单列表成功",
		"data": gin.H{
			"orders": orders,
		},
	})
}

// CheckServiceAccess 检查用户是否可以使用服务
func CheckServiceAccess(c *gin.Context) {
	// 获取当前用户ID
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 检查用户是否可以使用服务
	canAccess, err := model.CheckServiceAccess(userID)
	if err != nil {
		zap.L().Error("检查服务访问权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "检查服务访问权限失败",
			"data": nil,
		})
		return
	}

	if !canAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  "您今日的免费使用次数已用完，请购买会员卡或明天再来",
			"data": nil,
		})
		return
	}

	// 记录用户使用服务
	if err := model.UseService(userID); err != nil {
		zap.L().Error("记录服务使用失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "记录服务使用失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "可以使用服务",
		"data": nil,
	})
}
