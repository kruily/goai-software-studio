// Package taskqueue 定义异步任务队列的业务抽象。
//
// 设计原则（工作室规范）：
//   - 接口定义在本文件；具体实现放同包的 asynq.go / kafka.go 等。
//   - 业务层只表达 Enqueue/Consume 语义，不直接依赖 Asynq/Kafka/Pulsar SDK。
//   - 队列消息只携带 TaskID，任务详情回表读取，避免队列变成业务数据源。
//   - 新增队列实现请用 add-infra-adapter 技能；注册消费者请用 add-worker-task 技能。
//
// 调用链：api logic -> task 表落库 -> TaskQueue.Enqueue -> 消费者 Consume -> 回表执行。
package taskqueue

import (
	"context"
	"time"
)

// MatchMode 表示订阅与分发的匹配方式。
type MatchMode string

const (
	MatchExact  MatchMode = "exact"  // 精确匹配 topic。
	MatchPrefix MatchMode = "prefix" // 前缀匹配 topic。
	MatchGlob   MatchMode = "glob"   // 通配符匹配，便于映射 Kafka/Pulsar pattern。
	MatchRegex  MatchMode = "regex"  // 正则匹配。
)

// Config 收口队列适配器的通用配置。
type Config struct {
	Driver   string `json:",default=asynq"` // Driver 决定使用哪个适配器：asynq / kafka / pulsar。
	Queue    string `json:",optional"`      // Queue 是默认队列名。
	Group    string `json:",optional"`      // Group 是消费组名。
	Redis    string `json:",optional"`      // Redis 是 asynq 后端地址（asynq 专用）。
}

// Subscription 描述消费者想订阅的一组 topic。
type Subscription struct {
	Pattern string    // Pattern 是匹配表达式，例如 task.user. 或 ^task\.video\..+$。
	Mode    MatchMode // Mode 为空时按 exact 处理。
}

// EnqueueOptions 描述投递时的通用控制参数。
type EnqueueOptions struct {
	Queue     string        // Queue 为空时使用默认队列。
	MaxRetry  int           // MaxRetry 为 0 时使用实现默认值。
	Delay     time.Duration // Delay 表示延迟投递。
	Timeout   time.Duration // Timeout 限制单次处理耗时。
	Deadline  time.Time     // Deadline 表示最晚完成时间。
	UniqueFor time.Duration // UniqueFor 用于短时间幂等投递。
}

// Message 是消费者从队列拿到的统一消息。仅携带 TaskID，详情回表。
type Message struct {
	Topic  string
	TaskID int64
}

// Handler 定义处理任务消息的统一入口。
type Handler interface {
	Handle(ctx context.Context, msg Message) error
}

// HandlerFunc 允许普通函数作为 Handler。
type HandlerFunc func(ctx context.Context, msg Message) error

// Handle 执行任务处理函数。
func (f HandlerFunc) Handle(ctx context.Context, msg Message) error { return f(ctx, msg) }

// TaskQueue 是业务代码唯一依赖的异步任务队列抽象。
type TaskQueue interface {
	Enqueue(ctx context.Context, topic string, taskID int64, opts EnqueueOptions) error
	Consume(ctx context.Context, subscriptions []Subscription, group string, handler Handler) error
	Ack(ctx context.Context, msg Message) error
	Nack(ctx context.Context, msg Message, delay time.Duration) error
	Shutdown(ctx context.Context) error
}

// Payload 是队列里实际序列化的最小载荷。
type Payload struct {
	TaskID int64 `json:"task_id"`
}
