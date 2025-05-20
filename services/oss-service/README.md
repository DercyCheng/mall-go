# OSS Service

基于MongoDB GridFS的对象存储服务，为mall-go项目提供文件存储功能。

## 功能特性

- 文件上传/下载
- 文件元数据管理
- 文件类型验证
- 文件大小限制
- 多租户支持(通过ownerID)

## 技术栈

- Go 1.24+
- Gin Web框架
- MongoDB GridFS
- Viper配置管理
- Zap日志

## 快速开始

1. 启动MongoDB服务
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

2. 配置服务
修改`configs/config.yaml`文件:
```yaml
server:
  port: "8083"

mongodb:
  uri: "mongodb://localhost:27017"
  database: "oss_db"
```

3. 启动服务
```bash
go run cmd/main.go
```

## API文档

### 文件上传
```
POST /api/files/upload

参数:
- file: 上传的文件(表单字段)

响应:
{
  "file_id": "507f1f77bcf86cd799439011",
  "message": "file uploaded successfully"
}
```

### 文件下载
```
GET /api/files/:id
```

### 文件删除
```
DELETE /api/files/:id
```

### 获取文件信息
```
GET /api/files/:id/info

响应:
{
  "id": "507f1f77bcf86cd799439011",
  "filename": "example.jpg",
  "size": 1024,
  "mime_type": "image/jpeg",
  "created_at": "2023-01-01T00:00:00Z",
  "owner_id": "user123"
}
```