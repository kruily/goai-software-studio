// Package storage 定义对象存储的业务抽象。
//
// 设计原则（工作室规范）：
//   - 接口定义在本文件；具体厂商实现放同包的 minio.go / s3.go / oss.go 等。
//   - 业务层（logic/mq/cron）只依赖 Storage 接口，不直接依赖 MinIO/OSS/COS/S3 SDK，
//     避免后续换云厂商时迁移成本扩散到业务代码。
//   - 新增厂商适配器请用 add-infra-adapter 技能。
//
// 调用链：logic/worker -> Storage -> 具体厂商适配器。
package storage

import (
	"context"
	"io"
	"time"
)

// Config 收口对象存储适配器所需的通用配置。
type Config struct {
	Driver        string `json:",default=minio"` // Driver 决定使用哪个适配器：minio / s3 / oss / cos。
	Endpoint      string // Endpoint 是对象存储 API 地址，生产来自安全配置。
	AccessKeyID   string // AccessKeyID 是访问标识，不下发到客户端。
	AccessSecret  string // AccessSecret 是访问密钥，禁止写日志或返回给用户。
	UseSSL        bool   // UseSSL 控制是否走 HTTPS。
	Bucket        string // Bucket 是默认存储桶。
	Region        string // Region 记录存储区域。
	PublicBaseURL string // PublicBaseURL 是无 CDN 时的公开访问前缀。
	CDNBaseURL    string // CDNBaseURL 是生产优先使用的 CDN 前缀。
	CreateBucket  bool   // CreateBucket 仅建议开发环境开启，生产由基础设施提前创建。
}

// PresignOptions 描述预签名 URL 的业务约束。
type PresignOptions struct {
	Expires     time.Duration     // Expires 是签名有效期，0 时由适配器用安全默认值。
	ContentType string            // ContentType 供上传方设置 MIME 类型。
	Headers     map[string]string // Headers 随请求携带，禁止放服务端密钥。
}

// PresignResult 是返回给客户端的临时访问结果，仅在 ExpiresAt 前有效。
type PresignResult struct {
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers,omitempty"`
	ExpiresAt time.Time         `json:"expires_at"`
	ObjectKey string            `json:"object_key"`
}

// Storage 定义业务层可使用的对象存储能力。
// 业务代码不得直接依赖 MinIO/OSS/COS/S3 SDK。
type Storage interface {
	PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, key string) error
	PresignPut(ctx context.Context, key string, opts PresignOptions) (*PresignResult, error)
	PresignGet(ctx context.Context, key string, opts PresignOptions) (*PresignResult, error)
	PublicURL(ctx context.Context, key string) string
}

// DefaultPresignExpiry 是未显式指定有效期时的安全默认值，供各适配器复用。
func DefaultPresignExpiry(opts PresignOptions) time.Duration {
	if opts.Expires > 0 {
		return opts.Expires
	}
	return 15 * time.Minute
}
