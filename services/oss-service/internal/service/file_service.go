package service

import (
	"context"
	"errors"
	"io"
	"mall-go/services/oss-service/configs"
	"mall-go/services/oss-service/internal/repository"
	"mime"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file size exceeds limit")
)

type FileService struct {
	repo   repository.FileRepository
	cfg    *configs.StorageConfig
	logger *zap.Logger
}

func NewFileService(repo repository.FileRepository, cfg *configs.StorageConfig, logger *zap.Logger) *FileService {
	return &FileService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *FileService) UploadFile(ctx context.Context, filename string, file io.Reader, size int64, ownerID string) (string, error) {
	// 验证文件大小
	if size > s.cfg.MaxUploadSize {
		return "", ErrFileTooLarge
	}

	// 验证文件类型
	mimeType := mime.TypeByExtension(filename)
	if !s.isAllowedType(mimeType) {
		return "", ErrInvalidFileType
	}

	// 创建文件元数据
	metadata := &repository.FileMetadata{
		ID:        primitive.NewObjectID(),
		Filename:  filename,
		MimeType:  mimeType,
		CreatedAt: time.Now(),
		OwnerID:   ownerID,
	}

	// 上传文件
	fileID, err := s.repo.Upload(ctx, filename, file, metadata)
	if err != nil {
		s.logger.Error("Failed to upload file",
			zap.String("filename", filename),
			zap.Error(err))
		return "", err
	}

	return fileID, nil
}

func (s *FileService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, *repository.FileMetadata, error) {
	return s.repo.Download(ctx, fileID)
}

func (s *FileService) DeleteFile(ctx context.Context, fileID string) error {
	return s.repo.Delete(ctx, fileID)
}

func (s *FileService) GetFileInfo(ctx context.Context, fileID string) (*repository.FileMetadata, error) {
	return s.repo.GetMetadata(ctx, fileID)
}

func (s *FileService) isAllowedType(mimeType string) bool {
	for _, allowed := range s.cfg.AllowedTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

func (s *FileService) GetMaxUploadSize() int64 {
	return s.cfg.MaxUploadSize
}
