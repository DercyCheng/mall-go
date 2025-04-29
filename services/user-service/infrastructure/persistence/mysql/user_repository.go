package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

// UserEntity is the MySQL database entity for users
type UserEntity struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Username  string    `gorm:"column:username;uniqueIndex;not null;size:64"`
	Password  string    `gorm:"column:password;not null;size:128"`
	Email     string    `gorm:"column:email;uniqueIndex;size:100"`
	NickName  string    `gorm:"column:nick_name;size:200"`
	Icon      string    `gorm:"column:icon;size:500"`
	Status    int       `gorm:"column:status;default:1"`
	Note      string    `gorm:"column:note;size:500"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	LastLogin time.Time `gorm:"column:last_login"`
}

// TableName returns the table name for GORM
func (UserEntity) TableName() string {
	return "ums_admin"
}

// RoleEntity is the MySQL database entity for roles
type RoleEntity struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Name        string    `gorm:"column:name;uniqueIndex;not null;size:64"`
	Description string    `gorm:"column:description;size:500"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the table name for GORM
func (RoleEntity) TableName() string {
	return "ums_role"
}

// UserRoleEntity represents the many-to-many relationship between users and roles
type UserRoleEntity struct {
	ID     string `gorm:"column:id;primaryKey"`
	UserID string `gorm:"column:user_id;index;not null"`
	RoleID string `gorm:"column:role_id;index;not null"`
}

// TableName returns the table name for GORM
func (UserRoleEntity) TableName() string {
	return "ums_admin_role_relation"
}

// UserRepositoryImpl implements the UserRepository interface for MySQL
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new MySQL repository for users
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save persists a user to the database
func (r *UserRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	userEntity := mapUserToEntity(user)
	return r.db.WithContext(ctx).Create(userEntity).Error
}

// FindByID finds a user by their unique ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	var userEntity UserEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	user := mapEntityToUser(&userEntity)

	// Get user roles
	if err := r.loadUserRoles(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByUsername finds a user by their username
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var userEntity UserEntity
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	user := mapEntityToUser(&userEntity)

	// Get user roles
	if err := r.loadUserRoles(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail finds a user by their email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var userEntity UserEntity
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	user := mapEntityToUser(&userEntity)

	// Get user roles
	if err := r.loadUserRoles(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user's information
func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	userEntity := mapUserToEntity(user)
	return r.db.WithContext(ctx).Where("id = ?", user.ID).Updates(userEntity).Error
}

// Delete removes a user from the database
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	// First delete user role relations
	if err := r.db.WithContext(ctx).Where("user_id = ?", id).Delete(&UserRoleEntity{}).Error; err != nil {
		return err
	}

	// Then delete the user
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserEntity{}).Error
}

// List returns a list of users with pagination
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var userEntities []UserEntity
	var total int64

	// Count total users
	if err := r.db.WithContext(ctx).Model(&UserEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated users
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&userEntities).Error; err != nil {
		return nil, 0, err
	}

	// Map to domain model
	users := make([]*model.User, len(userEntities))
	for i, entity := range userEntities {
		users[i] = mapEntityToUser(&entity)
		if err := r.loadUserRoles(ctx, users[i]); err != nil {
			return nil, 0, err
		}
	}

	return users, total, nil
}

// Search searches for users based on criteria
func (r *UserRepositoryImpl) Search(ctx context.Context, query string, page, pageSize int) ([]*model.User, int64, error) {
	var userEntities []UserEntity
	var total int64

	// Build search query
	searchQuery := "%" + query + "%"
	dbQuery := r.db.WithContext(ctx).Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", searchQuery, searchQuery, searchQuery)

	// Count total matching users
	if err := dbQuery.Model(&UserEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&userEntities).Error; err != nil {
		return nil, 0, err
	}

	// Map to domain model
	users := make([]*model.User, len(userEntities))
	for i, entity := range userEntities {
		users[i] = mapEntityToUser(&entity)
		if err := r.loadUserRoles(ctx, users[i]); err != nil {
			return nil, 0, err
		}
	}

	return users, total, nil
}

// loadUserRoles loads the roles for a user
func (r *UserRepositoryImpl) loadUserRoles(ctx context.Context, user *model.User) error {
	var roleEntities []RoleEntity

	// Join query to get roles for user
	if err := r.db.WithContext(ctx).
		Joins("JOIN ums_admin_role_relation ur ON ur.role_id = ums_role.id").
		Where("ur.user_id = ?", user.ID).
		Find(&roleEntities).Error; err != nil {
		return err
	}

	// Map role entities to domain models
	roles := make([]model.Role, len(roleEntities))
	for i, entity := range roleEntities {
		rolePtr := mapEntityToRole(&entity)
		roles[i] = *rolePtr // Dereference to get the value
	}

	user.Roles = roles
	return nil
}

// RoleRepositoryImpl implements the RoleRepository interface for MySQL
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository creates a new MySQL repository for roles
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Save persists a role to the database
func (r *RoleRepositoryImpl) Save(ctx context.Context, role *model.Role) error {
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	roleEntity := mapRoleToEntity(role)
	return r.db.WithContext(ctx).Create(roleEntity).Error
}

// FindByID finds a role by its unique ID
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Role, error) {
	var roleEntity RoleEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&roleEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	role := mapEntityToRole(&roleEntity)
	return role, nil
}

// FindByName finds a role by its name
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var roleEntity RoleEntity
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&roleEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	role := mapEntityToRole(&roleEntity)
	return role, nil
}

// Update updates a role's information
func (r *RoleRepositoryImpl) Update(ctx context.Context, role *model.Role) error {
	roleEntity := mapRoleToEntity(role)
	return r.db.WithContext(ctx).Where("id = ?", role.ID).Updates(roleEntity).Error
}

// Delete removes a role from the database
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id string) error {
	// First delete user-role relations
	if err := r.db.WithContext(ctx).Where("role_id = ?", id).Delete(&UserRoleEntity{}).Error; err != nil {
		return err
	}

	// Then delete the role
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&RoleEntity{}).Error
}

// List returns a list of roles
func (r *RoleRepositoryImpl) List(ctx context.Context) ([]*model.Role, error) {
	var roleEntities []RoleEntity

	if err := r.db.WithContext(ctx).Find(&roleEntities).Error; err != nil {
		return nil, err
	}

	// Map to domain model
	roles := make([]*model.Role, len(roleEntities))
	for i, entity := range roleEntities {
		roles[i] = mapEntityToRole(&entity)
	}

	return roles, nil
}

// GetUserRoles returns all roles for a given user
func (r *RoleRepositoryImpl) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	var roleEntities []RoleEntity

	// Join query to get roles for user
	if err := r.db.WithContext(ctx).
		Joins("JOIN ums_admin_role_relation ur ON ur.role_id = ums_role.id").
		Where("ur.user_id = ?", userID).
		Find(&roleEntities).Error; err != nil {
		return nil, err
	}

	// Map to domain model
	roles := make([]*model.Role, len(roleEntities))
	for i, entity := range roleEntities {
		roles[i] = mapEntityToRole(&entity)
	}

	return roles, nil
}

// AssignRoleToUser assigns a role to a user
func (r *RoleRepositoryImpl) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	userRoleEntity := UserRoleEntity{
		ID:     uuid.New().String(),
		UserID: userID,
		RoleID: roleID,
	}

	return r.db.WithContext(ctx).Create(&userRoleEntity).Error
}

// RevokeRoleFromUser removes a role from a user
func (r *RoleRepositoryImpl) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRoleEntity{}).Error
}

// Helper functions to map between domain models and database entities

func mapUserToEntity(user *model.User) *UserEntity {
	return &UserEntity{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		NickName:  user.NickName,
		Icon:      user.Icon,
		Status:    int(user.Status),
		Note:      user.Note,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		LastLogin: user.LastLogin,
	}
}

func mapEntityToUser(entity *UserEntity) *model.User {
	return &model.User{
		ID:        entity.ID,
		Username:  entity.Username,
		Password:  entity.Password,
		Email:     entity.Email,
		NickName:  entity.NickName,
		Icon:      entity.Icon,
		Status:    model.UserStatus(entity.Status),
		Note:      entity.Note,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		LastLogin: entity.LastLogin,
		Roles:     make([]model.Role, 0),
	}
}

func mapRoleToEntity(role *model.Role) *RoleEntity {
	return &RoleEntity{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func mapEntityToRole(entity *RoleEntity) *model.Role {
	return &model.Role{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
