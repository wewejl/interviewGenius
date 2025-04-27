package model

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID            uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	Username      string     `json:"username" gorm:"size:64;not null;unique"`
	Password      string     `json:"-" gorm:"size:100;not null"`
	Email         string     `json:"email" gorm:"size:100;unique"`
	MemberExpiry  *time.Time `json:"member_expiry"`                    // 会员到期时间
	LastUseDate   *time.Time `json:"last_use_date"`                    // 记录未买卡用户最后使用日期
	DailyUseCount int        `json:"daily_use_count" gorm:"default:0"` // 当天已用次数
	Roles         []*Role    `json:"roles,omitempty" gorm:"many2many:user_role;"`
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeCreate 创建前生成UUID和加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if len(u.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// BeforeUpdate 更新前加密密码
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	if err := DB.Preload("Roles").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := DB.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// UpdateUser 更新用户
func UpdateUser(user *User) error {
	return DB.Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(id uuid.UUID) error {
	return DB.Delete(&User{}, "id = ?", id).Error
}

// GetUserList 获取用户列表
func GetUserList(page, limit int) ([]*User, int64, error) {
	var users []*User
	var total int64

	offset := (page - 1) * limit

	if err := DB.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AddRolesToUser 为用户添加角色
func AddRolesToUser(userID uuid.UUID, roleIDs []uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			return err
		}

		var roles []*Role
		if err := tx.Find(&roles, roleIDs).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Roles").Append(roles); err != nil {
			return err
		}

		return nil
	})
}

// GetUserRoles 获取用户的角色
func GetUserRoles(userID uuid.UUID) ([]*Role, error) {
	var user User
	if err := DB.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return user.Roles, nil
}

// RemoveRoleFromUser 从用户移除角色
func RemoveRoleFromUser(userID uuid.UUID, roleID uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			return err
		}

		var role Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Roles").Delete(&role); err != nil {
			return err
		}

		return nil
	})
}

// CheckServiceAccess 检查用户是否可以使用服务
func CheckServiceAccess(userID uuid.UUID) (bool, error) {
	var user User
	if err := DB.First(&user, "id = ?", userID).Error; err != nil {
		return false, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 检查是否是会员
	isMember := user.MemberExpiry != nil && user.MemberExpiry.After(now)

	// 如果是会员，直接允许访问
	if isMember {
		return true, nil
	}

	// 非会员用户检查每日使用次数
	// 检查是否是新的一天
	if user.LastUseDate == nil || user.LastUseDate.Before(today) {
		// 如果是新的一天，重置使用次数
		if err := DB.Model(&user).Updates(map[string]interface{}{
			"last_use_date":   today,
			"daily_use_count": 0,
		}).Error; err != nil {
			return false, err
		}
		user.DailyUseCount = 0
	}

	// 检查是否超过每日使用限制
	if user.DailyUseCount >= 1 {
		return false, nil
	}

	return true, nil
}

// UseService 记录用户使用服务
func UseService(userID uuid.UUID) error {
	var user User
	if err := DB.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 检查是否是会员
	isMember := user.MemberExpiry != nil && user.MemberExpiry.After(now)

	// 会员用户不消耗使用次数
	if isMember {
		return nil
	}

	// 检查是否是新的一天
	if user.LastUseDate == nil || user.LastUseDate.Before(today) {
		// 如果是新的一天，重置使用次数并设置为1
		return DB.Model(&user).Updates(map[string]interface{}{
			"last_use_date":   today,
			"daily_use_count": 1,
		}).Error
	} else {
		// 增加使用次数
		return DB.Model(&user).Updates(map[string]interface{}{
			"daily_use_count": user.DailyUseCount + 1,
		}).Error
	}
}

// IsMember 检查用户是否是会员
func IsMember(userID uuid.UUID) (bool, *time.Time, error) {
	var user User
	if err := DB.First(&user, "id = ?", userID).Error; err != nil {
		return false, nil, err
	}

	now := time.Now()
	isMember := user.MemberExpiry != nil && user.MemberExpiry.After(now)

	return isMember, user.MemberExpiry, nil
}
