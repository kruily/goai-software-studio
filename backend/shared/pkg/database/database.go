// Package database 提供 GORM 数据库连接的构造与健康检查，支持多种数据库。
//
// 设计：
//   - 数据库选型由项目决定（PostgreSQL / MySQL / SQLite 等），不写死。
//   - 由 Config.Driver 分发到对应 gorm driver；业务与 model 层只拿到 *gorm.DB，
//     不感知底层是哪种数据库。
//   - 具体驱动的 Open 实现由 bootstrap-project 按选型填入（见 openPostgres 等 TODO）。
//
// 调用链：config.DB -> database.New/MustOpen -> *gorm.DB -> ServiceContext.DB -> model。
package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Driver 表示数据库类型。
type Driver string

const (
	DriverPostgres Driver = "postgres" // PostgreSQL。
	DriverMySQL    Driver = "mysql"    // MySQL / MariaDB。
	DriverSQLite   Driver = "sqlite"   // SQLite，适合本地开发或轻量场景。
)

// Config 收口数据库连接配置，兼容多种数据库。
// DSN 相关字段按所选 Driver 取用；未用到的字段留空即可。
type Config struct {
	Driver          string `json:",default=postgres"` // Driver 决定使用哪种数据库：postgres / mysql / sqlite。
	Host            string `json:",optional"`
	Port            int    `json:",optional"`
	User            string `json:",optional"`
	Password        string `json:",optional"` // Password 来自安全配置，禁止入日志。
	DBName          string `json:",optional"`
	Params          string `json:",optional"` // Params 是附加连接参数（如 sslmode、charset、parseTime）。
	DSN             string `json:",optional"` // DSN 若显式提供则优先使用，忽略上面的分字段。
	MaxOpenConns    int    `json:",default=50"`
	MaxIdleConns    int    `json:",default=10"`
	ConnMaxLifetime int    `json:",default=3600"` // 单位秒。
}

// BuildDSN 依据 Driver 拼装连接串（DSN 显式提供时直接返回）。
func (c Config) BuildDSN() string {
	if c.DSN != "" {
		return c.DSN
	}
	switch Driver(c.Driver) {
	case DriverPostgres:
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
			c.Host, c.Port, c.User, c.Password, c.DBName, c.Params)
	case DriverMySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			c.User, c.Password, c.Host, c.Port, c.DBName, c.Params)
	case DriverSQLite:
		return c.DBName
	default:
		return c.DSN
	}
}

// New 是数据库连接的工厂入口，按 Config.Driver 分发。
func New(c Config) (*gorm.DB, error) {
	switch Driver(c.Driver) {
	case DriverPostgres:
		return openPostgres(c)
	case DriverMySQL:
		return openMySQL(c)
	case DriverSQLite:
		return openSQLite(c)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", c.Driver)
	}
}

// MustOpen 打开 GORM 连接，失败即 panic（进程启动阶段调用）。
func MustOpen(c Config) *gorm.DB {
	db, err := New(c)
	if err != nil {
		panic(fmt.Sprintf("database.MustOpen failed: %v", err))
	}
	return db
}

func openPostgres(c Config) (*gorm.DB, error) {
	// TODO(bootstrap): 引入 gorm.io/driver/postgres 后实现：
	//   dsn := c.BuildDSN()
	//   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//   sqlDB, _ := db.DB()
	//   sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	//   sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	//   sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)
	//   return db, nil
	return nil, fmt.Errorf("postgres driver not wired: add gorm.io/driver/postgres to go.mod and implement openPostgres")
}

func openMySQL(c Config) (*gorm.DB, error) {
	// TODO(bootstrap): 引入 gorm.io/driver/mysql 后实现。
	return nil, fmt.Errorf("mysql driver not wired: add gorm.io/driver/mysql to go.mod and implement openMySQL")
}

func openSQLite(c Config) (*gorm.DB, error) {
	// TODO(bootstrap): 引入 gorm.io/driver/sqlite 后实现。
	return nil, fmt.Errorf("sqlite driver not wired: add gorm.io/driver/sqlite to go.mod and implement openSQLite")
}
