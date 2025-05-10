package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"mall-go/services/cart-service/domain/model"
	"mall-go/services/cart-service/domain/repository"
	"mall-go/services/cart-service/infrastructure/config"
)

// MySQLRepository implements the CartRepository interface using MySQL
type MySQLRepository struct {
	db *sqlx.DB
}

// NewMySQLRepository creates a new MySQL repository
func NewMySQLRepository(config config.DatabaseConfig) (repository.CartRepository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	repo := &MySQLRepository{
		db: db,
	}

	// Ensure tables exist
	if err := repo.createTablesIfNotExist(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

// createTablesIfNotExist creates the cart_items table if it doesn't exist
func (r *MySQLRepository) createTablesIfNotExist() error {
	query := `
	CREATE TABLE IF NOT EXISTS cart_items (
		id VARCHAR(50) PRIMARY KEY,
		user_id VARCHAR(50) NOT NULL,
		product_id VARCHAR(50) NOT NULL,
		product_sku_id VARCHAR(50) NOT NULL,
		product_name VARCHAR(255) NOT NULL,
		product_sub_title VARCHAR(255),
		product_pic VARCHAR(1000),
		product_price DECIMAL(10,2) NOT NULL,
		product_quantity INT NOT NULL,
		product_sku_code VARCHAR(50),
		product_category_id VARCHAR(50),
		product_brand VARCHAR(100),
		product_sn VARCHAR(100),
		product_attr TEXT,
		promotion_info VARCHAR(500),
		promotion_amount DECIMAL(10,2) DEFAULT 0,
		coupon_amount DECIMAL(10,2) DEFAULT 0,
		real_amount DECIMAL(10,2) NOT NULL,
		check_status BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_user_id (user_id),
		INDEX idx_product_id (product_id),
		UNIQUE INDEX idx_user_product_sku (user_id, product_id, product_sku_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := r.db.Exec(query)
	return err
}

// AddItem adds a new item to the cart
func (r *MySQLRepository) AddItem(ctx context.Context, item *model.CartItem) error {
	query := `
	INSERT INTO cart_items (
		id, user_id, product_id, product_sku_id, product_name, product_sub_title,
		product_pic, product_price, product_quantity, product_sku_code,
		product_category_id, product_brand, product_sn, product_attr,
		promotion_info, promotion_amount, coupon_amount, real_amount, check_status,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.UserID, item.ProductID, item.ProductSKUID, item.ProductName,
		item.ProductSubTitle, item.ProductPic, item.ProductPrice, item.ProductQuantity,
		item.ProductSkuCode, item.ProductCategoryID, item.ProductBrand, item.ProductSn,
		item.ProductAttr, item.PromotionInfo, item.PromotionAmount, item.CouponAmount,
		item.RealAmount, item.CheckStatus, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

// GetItem retrieves a cart item by ID
func (r *MySQLRepository) GetItem(ctx context.Context, id string) (*model.CartItem, error) {
	query := `SELECT * FROM cart_items WHERE id = ?`
	var item model.CartItem
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart item not found with ID: %s", id)
		}
		return nil, err
	}
	return &item, nil
}

// GetItemsByUserID retrieves all cart items for a user
func (r *MySQLRepository) GetItemsByUserID(ctx context.Context, userID string) ([]*model.CartItem, error) {
	query := `SELECT * FROM cart_items WHERE user_id = ? ORDER BY updated_at DESC`
	items := []*model.CartItem{}
	err := r.db.SelectContext(ctx, &items, query, userID)
	return items, err
}

// GetCheckedItemsByUserID retrieves all checked cart items for a user
func (r *MySQLRepository) GetCheckedItemsByUserID(ctx context.Context, userID string) ([]*model.CartItem, error) {
	query := `SELECT * FROM cart_items WHERE user_id = ? AND check_status = true ORDER BY updated_at DESC`
	items := []*model.CartItem{}
	err := r.db.SelectContext(ctx, &items, query, userID)
	return items, err
}

// UpdateItem updates an existing cart item
func (r *MySQLRepository) UpdateItem(ctx context.Context, item *model.CartItem) error {
	query := `
	UPDATE cart_items SET
		product_quantity = ?,
		product_price = ?,
		promotion_info = ?,
		promotion_amount = ?,
		coupon_amount = ?,
		real_amount = ?,
		check_status = ?,
		updated_at = ?
	WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		item.ProductQuantity, item.ProductPrice, item.PromotionInfo,
		item.PromotionAmount, item.CouponAmount, item.RealAmount,
		item.CheckStatus, time.Now(), item.ID,
	)
	return err
}

// DeleteItem deletes a cart item by ID
func (r *MySQLRepository) DeleteItem(ctx context.Context, id string) error {
	query := `DELETE FROM cart_items WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteItemsByUserID deletes all cart items for a user
func (r *MySQLRepository) DeleteItemsByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM cart_items WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// DeleteCheckedItemsByUserID deletes all checked cart items for a user
func (r *MySQLRepository) DeleteCheckedItemsByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM cart_items WHERE user_id = ? AND check_status = true`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// GetItemCount returns the number of cart items for a user
func (r *MySQLRepository) GetItemCount(ctx context.Context, userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM cart_items WHERE user_id = ?`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}

// GetCartTotalAmount returns the total amount of all checked items in a user's cart
func (r *MySQLRepository) GetCartTotalAmount(ctx context.Context, userID string) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(real_amount), 0) FROM cart_items WHERE user_id = ? AND check_status = true`
	err := r.db.GetContext(ctx, &total, query, userID)
	return total, err
}

// CheckExistingItem checks if an item already exists in the cart
func (r *MySQLRepository) CheckExistingItem(ctx context.Context, userID, productID, skuID string) (*model.CartItem, error) {
	query := `SELECT * FROM cart_items WHERE user_id = ? AND product_id = ? AND product_sku_id = ?`
	var item model.CartItem
	err := r.db.GetContext(ctx, &item, query, userID, productID, skuID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// UpdateCheckStatus updates the check status of a cart item
func (r *MySQLRepository) UpdateCheckStatus(ctx context.Context, id string, status bool) error {
	query := `UPDATE cart_items SET check_status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// UpdateAllCheckStatus updates the check status of all cart items for a user
func (r *MySQLRepository) UpdateAllCheckStatus(ctx context.Context, userID string, status bool) error {
	query := `UPDATE cart_items SET check_status = ?, updated_at = ? WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), userID)
	return err
}

// Clear clears all cart items in the database (for testing)
func (r *MySQLRepository) Clear(ctx context.Context) error {
	query := `DELETE FROM cart_items`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
