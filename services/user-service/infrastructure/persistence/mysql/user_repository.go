package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mall-go/pkg/database"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

// UserEntity 用户实体，用于ORM映射
type UserEntity struct {
	ID               string    `gorm:"primaryKey;type:varchar(36)"`
	Username         string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Password         string    `gorm:"type:varchar(255);not null"`
	Email            string    `gorm:"type:varchar(100);uniqueIndex"`
	Phone            string    `gorm:"type:varchar(20);uniqueIndex"`
	NickName         string    `gorm:"type:varchar(50)"`
	Icon             string    `gorm:"type:varchar(255)"`
	Gender           int       `gorm:"default:0"`
	Birthday         *time.Time
	City             string    `gorm:"type:varchar(50)"`
	Job              string    `gorm:"type:varchar(50)"`
	Status           string    `gorm:"type:varchar(10);default:'active'"`
	IntegrationPoints int       `gorm:"default:0"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
	LastLoginAt      *time.Time
}

// UserRoleRelation 用户角色关联实体
type UserRoleRelation struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"type:varchar(36);not null;index"`
	RoleID    string    `gorm:"type:varchar(36);not null;index"`
	CreatedAt time.Time `gorm:"not null"`
}

// TableName 设置表名
func (UserEntity) TableName() string {
	return "user"
}

func (UserRoleRelation) TableName() string {
	return "user_role_relation"
}

// UserRepositoryImpl 用户仓储MySQL实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{
		db: database.DB,
	}
}

// Save 保存用户
func (r *UserRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	// 将领域模型转换为持久化实体
	userEntity := &UserEntity{
		ID:               user.ID,
		Username:         user.Username,
		Password:         user.Password,
		Email:            user.Email,
		Phone:            user.Phone,
		NickName:         user.NickName,
		Icon:             user.Icon,
		Gender:           user.Gender,
		Birthday:         user.Birthday,
		City:             user.City,
		Job:              user.Job,
		Status:           string(user.Status),
		IntegrationPoints: user.IntegrationPoints,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		LastLoginAt:      user.LastLoginAt,
	}

	return r.db.Create(userEntity).Error
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	var userEntity UserEntity
	err := r.db.Where("id = ?", id).First(&userEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&userEntity)
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var userEntity UserEntity
	err := r.db.Where("username = ?", username).First(&userEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found by username: %s", username)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&userEntity)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var userEntity UserEntity
	err := r.db.Where("email = ?", email).First(&userEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found by email: %s", email)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&userEntity)
}

// FindByPhone 根据手机号查找用户
func (r *UserRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var userEntity UserEntity
	err := r.db.Where("phone = ?", phone).First(&userEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found by phone: %s", phone)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&userEntity)
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	// 将领域模型转换为持久化实体
	userEntity := &UserEntity{
		ID:               user.ID,
		Username:         user.Username,
		Password:         user.Password,
		Email:            user.Email,
		Phone:            user.Phone,
		NickName:         user.NickName,
		Icon:             user.Icon,
		Gender:           user.Gender,
		Birthday:         user.Birthday,
		City:             user.City,
		Job:              user.Job,
		Status:           string(user.Status),
		IntegrationPoints: user.IntegrationPoints,
		UpdatedAt:        user.UpdatedAt,
		LastLoginAt:      user.LastLoginAt,
	}

	return r.db.Model(&UserEntity{}).Where("id = ?", user.ID).Updates(userEntity).Error
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 开始事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Where("user_id = ?", id).Delete(&UserRoleRelation{}).Error; err != nil {
			return err
		}

		// 删除用户
		if err := tx.Where("id = ?", id).Delete(&UserEntity{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindAll 获取用户列表
func (r *UserRepositoryImpl) FindAll(ctx context.Context, page, size int) ([]*model.User, int64, error) {
	var userEntities []UserEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&UserEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Offset(offset).Limit(size).Order("created_at desc").Find(&userEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	users := make([]*model.User, 0, len(userEntities))
	for _, entity := range userEntities {
		user, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// AddUserRole 添加用户角色
func (r *UserRepositoryImpl) AddUserRole(ctx context.Context, userID, roleID string) error {
	relation := UserRoleRelation{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}

	return r.db.Create(&relation).Error
}

// RemoveUserRole 移除用户角色
func (r *UserRepositoryImpl) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRoleRelation{}).Error
}

// GetUserRoles 获取用户角色
func (r *UserRepositoryImpl) GetUserRoles(ctx context.Context, userID string) ([]model.Role, error) {
	// 查询用户角色关联
	var relations []UserRoleRelation
	if err := r.db.Where("user_id = ?", userID).Find(&relations).Error; err != nil {
		return nil, err
	}

	// 如果没有角色，返回空数组
	if len(relations) == 0 {
		return []model.Role{}, nil
	}

	// 提取角色ID
	roleIDs := make([]string, 0, len(relations))
	for _, relation := range relations {
		roleIDs = append(roleIDs, relation.RoleID)
	}

	// 查询角色
	var roleEntities []RoleEntity
	if err := r.db.Where("id IN ?", roleIDs).Find(&roleEntities).Error; err != nil {
		return nil, err
	}

	// 转换为领域模型
	roles := make([]model.Role, 0, len(roleEntities))
	for _, entity := range roleEntities {
		roles = append(roles, model.Role{
			ID:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
			Status:      model.RoleStatus(entity.Status),
			Sort:        entity.Sort,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		})
	}

	return roles, nil
}

// UpdateStatus 更新用户状态
func (r *UserRepositoryImpl) UpdateStatus(ctx context.Context, id string, status model.UserStatus) error {
	return r.db.Model(&UserEntity{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateLastLoginTime 更新用户登录时间
func (r *UserRepositoryImpl) UpdateLastLoginTime(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.Model(&UserEntity{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": now,
		"updated_at":    now,
	}).Error
}

// Search 搜索用户
func (r *UserRepositoryImpl) Search(ctx context.Context, keyword string, page, size int) ([]*model.User, int64, error) {
	var userEntities []UserEntity
	var total int64

	query := r.db.Model(&UserEntity{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nick_name LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at desc").Find(&userEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	users := make([]*model.User, 0, len(userEntities))
	for _, entity := range userEntities {
		user, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

// 辅助方法: 将持久化实体转换为领域模型
func (r *UserRepositoryImpl) toDomainModel(entity *UserEntity) (*model.User, error) {
	return &model.User{
		ID:               entity.ID,
		Username:         entity.Username,
		Password:         entity.Password,
		Email:            entity.Email,
		Phone:            entity.Phone,
		NickName:         entity.NickName,
		Icon:             entity.Icon,
		Gender:           entity.Gender,
		Birthday:         entity.Birthday,
		City:             entity.City,
		Job:              entity.Job,
		Status:           model.UserStatus(entity.Status),
		IntegrationPoints: entity.IntegrationPoints,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		LastLoginAt:      entity.LastLoginAt,
	}, nil
}