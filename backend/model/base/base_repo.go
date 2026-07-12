// Package base 提供 GORM 模型的通用 CRUD 能力。
//
// 设计（通用 BaseRepo 模式）：
//   - BaseRepo[T] 是泛型接口，定义最常见的 CRUD 操作。
//   - 业务 model 的 struct 直接声明通用字段（Id/UserId/CreatedBy/CreatedAt/UpdatedAt/DeletedAt），
//     而非嵌入一个 BaseModel struct——各域可按需调整。
//   - 业务 model 的 interface 嵌入 BaseRepo[T] 获得通用能力，再补充领域查询方法；
//     私有 struct 嵌入 baseRepo[T] 获得通用实现。
//
// 调用链：logic -> model repo -> BaseRepo 通用方法 -> gorm.DB -> 数据库。
package base

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BaseRepo 抽象出最常见的 GORM CRUD 操作。
type BaseRepo[T any] interface {
	Create(session *gorm.DB, data *T) error
	CreateBatch(session *gorm.DB, data []*T) error
	FirstOrCreate(session *gorm.DB, data *T, conds ...interface{}) error
	FindOneById(session *gorm.DB, id any) (*T, error)
	FindOneByIdForUpdate(session *gorm.DB, id any) (*T, error)
	FindList(session *gorm.DB, limit ...int) ([]*T, error)
	FindListForPage(session *gorm.DB, page int, pageSize int) ([]*T, error)
	Updates(session *gorm.DB, newData interface{}) error
	Update(session *gorm.DB, id any, newData interface{}) error
	Delete(session *gorm.DB, id any) error
	Count(session *gorm.DB) int64
}

// baseRepo 是 BaseRepo 的默认泛型实现。
type baseRepo[T any] struct{}

// NewBaseRepo 返回通用 repo 实例。
func NewBaseRepo[T any]() BaseRepo[T] {
	return &baseRepo[T]{}
}

func (b *baseRepo[T]) model(session *gorm.DB) *gorm.DB {
	var m T
	return session.Model(&m)
}

func (b *baseRepo[T]) Create(session *gorm.DB, data *T) error {
	return b.model(session).Create(data).Error
}

func (b *baseRepo[T]) CreateBatch(session *gorm.DB, data []*T) error {
	return b.model(session).CreateInBatches(data, len(data)).Error
}

func (b *baseRepo[T]) FirstOrCreate(session *gorm.DB, data *T, conds ...interface{}) error {
	return b.model(session).FirstOrCreate(data, conds...).Error
}

func (b *baseRepo[T]) FindOneById(session *gorm.DB, id any) (*T, error) {
	var resp T
	err := b.model(session).Last(&resp, "id = ?", id).Error
	return &resp, err
}

func (b *baseRepo[T]) FindOneByIdForUpdate(session *gorm.DB, id any) (*T, error) {
	var resp T
	err := b.model(session).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).Last(&resp).Error
	return &resp, err
}

func (b *baseRepo[T]) FindList(session *gorm.DB, limit ...int) ([]*T, error) {
	var models []*T
	sess := b.model(session)
	if len(limit) > 0 {
		sess = sess.Limit(limit[0])
	}
	err := sess.Find(&models).Error
	return models, err
}

func (b *baseRepo[T]) FindListForPage(session *gorm.DB, page int, pageSize int) ([]*T, error) {
	var models []*T
	err := b.model(session).Offset((page-1)*pageSize).Limit(pageSize).Find(&models).Error
	return models, err
}

func (b *baseRepo[T]) Updates(session *gorm.DB, newData interface{}) error {
	return b.model(session).Updates(newData).Error
}

func (b *baseRepo[T]) Update(session *gorm.DB, id any, newData interface{}) error {
	return b.model(session).Where("id = ?", id).Updates(newData).Error
}

func (b *baseRepo[T]) Delete(session *gorm.DB, id any) error {
	var m T
	return b.model(session).Where("id = ?", id).Delete(&m).Error
}

func (b *baseRepo[T]) Count(session *gorm.DB) int64 {
	var count int64
	b.model(session).Count(&count)
	return count
}
