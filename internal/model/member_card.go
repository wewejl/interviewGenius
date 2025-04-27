package model

// MemberCard 会员卡商品模型
type MemberCard struct {
	Model
	Name         string `json:"name" gorm:"size:20;not null"`         // 日卡/周卡等
	DurationDays int    `json:"duration_days" gorm:"not null"`        // 有效天数（日卡=1，周卡=7...）
	Price        int    `json:"price" gorm:"not null"`                // 价格（单位：分）
	Description  string `json:"description" gorm:"size:255;not null"` // 会员卡描述
}

// GetAllMemberCards 获取所有会员卡类型
func GetAllMemberCards() ([]*MemberCard, error) {
	var cards []*MemberCard
	if err := DB.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

// GetMemberCardByID 根据ID获取会员卡
func GetMemberCardByID(id uint) (*MemberCard, error) {
	var card MemberCard
	if err := DB.First(&card, id).Error; err != nil {
		return nil, err
	}
	return &card, nil
}
