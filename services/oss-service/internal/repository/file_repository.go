package repository

import (
	"context"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.uber.org/zap"
)

type FileMetadata struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Filename  string             `bson:"filename"`
	Size      int64              `bson:"size"`
	MimeType  string             `bson:"mime_type"`
	CreatedAt time.Time          `bson:"created_at"`
	OwnerID   string             `bson:"owner_id,omitempty"`
}

type FileRepository interface {
	Upload(ctx context.Context, filename string, file io.Reader, metadata *FileMetadata) (string, error)
	Download(ctx context.Context, fileID string) (io.ReadCloser, *FileMetadata, error)
	Delete(ctx context.Context, fileID string) error
	GetMetadata(ctx context.Context, fileID string) (*FileMetadata, error)
}

type GridFSRepository struct {
	bucket *gridfs.Bucket
	logger *zap.Logger
}

func NewGridFSRepository(bucket *gridfs.Bucket, logger *zap.Logger) *GridFSRepository {
	return &GridFSRepository{
		bucket: bucket,
		logger: logger,
	}
}

func (r *GridFSRepository) Upload(ctx context.Context, filename string, file io.Reader, metadata *FileMetadata) (string, error) {
	uploadStream, err := r.bucket.OpenUploadStreamWithID(metadata.ID, filename)
	if err != nil {
		return "", err
	}
	defer uploadStream.Close()

	size, err := io.Copy(uploadStream, file)
	if err != nil {
		return "", err
	}

	metadata.Size = size
	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

func (r *GridFSRepository) Download(ctx context.Context, fileID string) (io.ReadCloser, *FileMetadata, error) {
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, nil, err
	}

	downloadStream, err := r.bucket.OpenDownloadStream(objectID)
	if err != nil {
		return nil, nil, err
	}
	fileInfo := downloadStream.GetFile()
	metadata := &FileMetadata{
		ID:        fileInfo.ID.(primitive.ObjectID),
		Filename:  fileInfo.Name,
		Size:      fileInfo.Length,
		CreatedAt: fileInfo.UploadDate,
	}

	// Decode metadata to extract mime_type
	if fileInfo.Metadata != nil {
		var metadataMap map[string]interface{}
		if err := bson.Unmarshal(fileInfo.Metadata, &metadataMap); err == nil {
			if mimeType, ok := metadataMap["mime_type"].(string); ok {
				metadata.MimeType = mimeType
			}
		}
	}
	return downloadStream, metadata, nil
}

func (r *GridFSRepository) Delete(ctx context.Context, fileID string) error {
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return err
	}
	return r.bucket.Delete(objectID)
}
func (r *GridFSRepository) GetMetadata(ctx context.Context, fileID string) (*FileMetadata, error) {
	objectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}

	var result FileMetadata
	err = r.bucket.GetFilesCollection().FindOne(ctx, primitive.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
