// Package custom_type 存放 GORM 结构化 JSON 字段的自定义类型。
//
// 示例：JSONStruct 是按数据库方言返回 jsonb / json / text 的泛型包装。
// 各域在 model/custom_type/{domain}_data.go 中定义自己的结构体并嵌入该模式。
package custom_type

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// JSONStruct 提供跨数据库兼容的 JSON 字段自定义类型。
// 嵌入到域内的结构体中使用，自动按方言映射到 jsonb(PG)/json(MySQL)/text(SQLite)。
type JSONStruct[T any] struct {
	Data T
}

// GormDBDataType 按数据库方言返回物理类型。
func (j JSONStruct[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "jsonb"
	case "mysql":
		return "json"
	default: // sqlite 等
		return "text"
	}
}

// Value 序列化为 JSON 字符串（实现 driver.Valuer）。
func (j JSONStruct[T]) Value() (driver.Value, error) {
	if j.isZero() {
		return nil, nil
	}
	return json.Marshal(j.Data)
}

// Scan 反序列化为结构体（实现 sql.Scanner）。
func (j *JSONStruct[T]) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, &j.Data)
	case string:
		return json.Unmarshal([]byte(v), &j.Data)
	default:
		return fmt.Errorf("unsupported scan type for JSONStruct: %T", src)
	}
}

func (j JSONStruct[T]) isZero() bool {
	var zero T
	return any(j.Data) == any(zero)
}
