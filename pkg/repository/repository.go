// Package repository provides common repository patterns and interfaces
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Common errors
var (
	ErrNotFound      = errors.New("entity not found")
	ErrInvalidEntity = errors.New("invalid entity")
	ErrDatabase      = errors.New("database error")
)

// Repository defines common CRUD operations for any entity
type Repository[T any, ID any] interface {
	// Create persists a new entity
	Create(ctx context.Context, entity *T) error

	// FindByID retrieves an entity by its ID
	FindByID(ctx context.Context, id ID) (*T, error)

	// Update updates an existing entity
	Update(ctx context.Context, entity *T) error

	// Delete removes an entity by its ID
	Delete(ctx context.Context, id ID) error

	// FindAll retrieves all entities
	FindAll(ctx context.Context) ([]T, error)
}

// GormRepository is a generic implementation of Repository using GORM
type GormRepository[T any, ID any] struct {
	db *gorm.DB
}

// NewGormRepository creates a new GormRepository
func NewGormRepository[T any, ID any](db *gorm.DB) *GormRepository[T, ID] {
	return &GormRepository[T, ID]{
		db: db,
	}
}

// Create persists a new entity
func (r *GormRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	if entity == nil {
		return ErrInvalidEntity
	}

	result := r.db.WithContext(ctx).Create(entity)
	if result.Error != nil {
		return wrapError(result.Error)
	}

	return nil
}

// FindByID retrieves an entity by its ID
func (r *GormRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	var entity T

	result := r.db.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, wrapError(result.Error)
	}

	return &entity, nil
}

// Update updates an existing entity
func (r *GormRepository[T, ID]) Update(ctx context.Context, entity *T) error {
	if entity == nil {
		return ErrInvalidEntity
	}

	result := r.db.WithContext(ctx).Save(entity)
	if result.Error != nil {
		return wrapError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete removes an entity by its ID
func (r *GormRepository[T, ID]) Delete(ctx context.Context, id ID) error {
	var entity T

	result := r.db.WithContext(ctx).Delete(&entity, id)
	if result.Error != nil {
		return wrapError(result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// FindAll retrieves all entities
func (r *GormRepository[T, ID]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T

	result := r.db.WithContext(ctx).Find(&entities)
	if result.Error != nil {
		return nil, wrapError(result.Error)
	}

	return entities, nil
}

// FindWithPagination retrieves entities with pagination
func (r *GormRepository[T, ID]) FindWithPagination(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, wrapError(err)
	}

	// Get records with pagination
	result := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities)
	if result.Error != nil {
		return nil, 0, wrapError(result.Error)
	}

	return entities, total, nil
}

// FindWithFilter retrieves entities with a filter
func (r *GormRepository[T, ID]) FindWithFilter(ctx context.Context, query interface{}, args ...interface{}) ([]T, error) {
	var entities []T

	result := r.db.WithContext(ctx).Where(query, args...).Find(&entities)
	if result.Error != nil {
		return nil, wrapError(result.Error)
	}

	return entities, nil
}

// FindWithFilterAndPagination retrieves entities with a filter and pagination
func (r *GormRepository[T, ID]) FindWithFilterAndPagination(
	ctx context.Context,
	query interface{},
	page, pageSize int,
	args ...interface{},
) ([]T, int64, error) {
	var entities []T
	var total int64

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count with filter
	countQuery := r.db.WithContext(ctx).Model(new(T)).Where(query, args...)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, wrapError(err)
	}

	// Get records with filter and pagination
	result := r.db.WithContext(ctx).Where(query, args...).Offset(offset).Limit(pageSize).Find(&entities)
	if result.Error != nil {
		return nil, 0, wrapError(result.Error)
	}

	return entities, total, nil
}

// wrapError wraps a GORM error into a repository error
func wrapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	// Here you could add more specific error handling
	return errors.Join(ErrDatabase, err)
}
