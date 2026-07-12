// Package model 是 GORM 模型目录。
//
// 迁移入口。按 model/{domain}/{name}.go 模式创建模型后在 AutoMigrate 追加。
package model

import (
	"gorm.io/gorm"
)

// Migrate 自动迁移注册。新增模型时在此追加。
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// &user.User{},
		// 新增模型在此追加...
	)
}
