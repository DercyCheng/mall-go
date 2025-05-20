package infrastructure

import (
	"context"
	"fmt"
	"mall-go/services/oss-service/configs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	GridFS   *gridfs.Bucket
}

func NewMongoDB(cfg *configs.MongoDBConfig, logger *zap.Logger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create GridFS bucket: %w", err)
	}

	logger.Info("Successfully connected to MongoDB",
		zap.String("database", cfg.Database))

	return &MongoDB{
		Client:   client,
		Database: db,
		GridFS:   bucket,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}
