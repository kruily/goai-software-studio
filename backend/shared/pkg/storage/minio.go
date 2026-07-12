// Package storage 提供对象存储的具体厂商实现。
//
// 当前实现：
//   - minio.go   — MinIO / S3 / OSS / COS / R2（统一通过 minio-go v7 SDK 对接）
//   - （其他厂商按需新增，见 add-infra-adapter 技能）
//
// 设计：每个实现都桥接到本包顶层的 Storage 接口，业务只依赖接口。
package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minioStore 通过 minio-go v7 SDK 实现 Storage 接口。
// MinIO 自托管、阿里云 OSS、腾讯云 COS、AWS S3、Cloudflare R2 均兼容此 SDK，
// 差异仅在于 Endpoint 与 Region，由 Config 控制。
type minioStore struct {
	client     *minio.Client
	bucket     string
	publicBase string
	cdnBase    string
}

// NewMinioStorage 构造 MinIO 兼容的存储适配器。
func NewMinioStorage(c Config) (Storage, error) {
	if c.Endpoint == "" {
		return nil, fmt.Errorf("storage endpoint cannot be empty")
	}
	if c.Bucket == "" {
		return nil, fmt.Errorf("storage bucket cannot be empty")
	}

	client, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.AccessSecret, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client failed: %w", err)
	}

	// 开发环境自动创建 bucket；生产环境应提前建好。
	if c.CreateBucket {
		ctx := context.Background()
		exists, err := client.BucketExists(ctx, c.Bucket)
		if err == nil && !exists {
			_ = client.MakeBucket(ctx, c.Bucket, minio.MakeBucketOptions{})
		}
	}

	return &minioStore{
		client:     client,
		bucket:     c.Bucket,
		publicBase: c.PublicBaseURL,
		cdnBase:    c.CDNBaseURL,
	}, nil
}

// PutObject 上传对象。
func (s *minioStore) PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("put object %s failed: %w", key, err)
	}
	return nil
}

// GetObject 获取对象流。
func (s *minioStore) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object %s failed: %w", key, err)
	}
	return obj, nil
}

// DeleteObject 删除对象。
func (s *minioStore) DeleteObject(ctx context.Context, key string) error {
	if err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object %s failed: %w", key, err)
	}
	return nil
}

// PresignPut 生成上传预签名 URL。
func (s *minioStore) PresignPut(ctx context.Context, key string, opts PresignOptions) (*PresignResult, error) {
	expiry := DefaultPresignExpiry(opts)
	u, err := s.client.PresignedPutObject(ctx, s.bucket, key, expiry)
	if err != nil {
		return nil, fmt.Errorf("presign put %s failed: %w", key, err)
	}
	return &PresignResult{
		URL:       u.String(),
		Method:    "PUT",
		Headers:   map[string]string{"Content-Type": opts.ContentType},
		ExpiresAt: time.Now().Add(expiry),
		ObjectKey: key,
	}, nil
}

// PresignGet 生成下载预签名 URL。
func (s *minioStore) PresignGet(ctx context.Context, key string, opts PresignOptions) (*PresignResult, error) {
	expiry := DefaultPresignExpiry(opts)
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, expiry, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("presign get %s failed: %w", key, err)
	}
	return &PresignResult{
		URL:       u.String(),
		Method:    "GET",
		ExpiresAt: time.Now().Add(expiry),
		ObjectKey: key,
	}, nil
}

// PublicURL 返回对象的公开访问 URL。
func (s *minioStore) PublicURL(ctx context.Context, key string) string {
	if s.cdnBase != "" {
		return s.cdnBase + "/" + key
	}
	if s.publicBase != "" {
		return s.publicBase + "/" + key
	}
	return ""
}
