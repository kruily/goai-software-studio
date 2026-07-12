// Package redisx 提供 Redis 客户端配置与构造。
// 供缓存、分布式锁、以及 asynq 队列后端复用。
//
// 调用链：config.Redis -> redisx.MustNew -> 客户端 -> 业务缓存 / taskqueue asynq 适配器。
package redisx

// Config 收口 Redis 连接配置。
type Config struct {
	Addr     string `json:",default=127.0.0.1:6379"`
	Password string `json:",optional"` // Password 来自安全配置，禁止入日志。
	DB       int    `json:",default=0"`
}
