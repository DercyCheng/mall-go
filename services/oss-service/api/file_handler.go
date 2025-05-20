package api

import (
	"errors"
	"io"
	"mall-go/services/oss-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type FileHandler struct {
	service *service.FileService
	logger  *zap.Logger
}

func NewFileHandler(service *service.FileService, logger *zap.Logger) *FileHandler {
	return &FileHandler{
		service: service,
		logger:  logger,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get uploaded file", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file upload"})
		return
	}
	defer file.Close()

	// 获取所有者ID(从JWT token中)
	ownerID := c.GetString("user_id")
	if ownerID == "" {
		ownerID = "anonymous"
	}

	// 调用服务上传文件
	fileID, err := h.service.UploadFile(c.Request.Context(), header.Filename, file, header.Size, ownerID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidFileType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
			return
		}
		if errors.Is(err, service.ErrFileTooLarge) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "file too large, max size is " +
					strconv.FormatInt(h.service.GetMaxUploadSize(), 10) + " bytes",
			})
			return
		}

		h.logger.Error("Failed to upload file",
			zap.String("filename", header.Filename),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id": fileID,
		"message": "file uploaded successfully",
	})
}

func (h *FileHandler) Delete(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file ID is required"})
		return
	}

	err := h.service.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		h.logger.Error("Failed to delete file",
			zap.String("fileID", fileID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted successfully"})
}

func (h *FileHandler) Download(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file ID is required"})
		return
	}

	file, metadata, err := h.service.DownloadFile(c.Request.Context(), fileID)
	if err != nil {
		h.logger.Error("Failed to download file",
			zap.String("fileID", fileID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download file"})
		return
	}
	defer file.Close()

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+metadata.Filename)
	c.Header("Content-Type", metadata.MimeType)
	c.Header("Content-Length", strconv.FormatInt(metadata.Size, 10))

	// 流式传输文件内容
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, file)
		return err == nil
	})
}

func (h *FileHandler) GetInfo(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file ID is required"})
		return
	}

	metadata, err := h.service.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		h.logger.Error("Failed to get file info",
			zap.String("fileID", fileID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file info"})
		return
	}

	c.JSON(http.StatusOK, metadata)
}
