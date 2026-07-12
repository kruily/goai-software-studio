// Package taskqueue 提供异步任务队列的具体实现。
//
// 当前实现：
//   - asynq.go — 基于 Asynq + Redis（工作室默认）
//   - （其他队列按需新增，见 add-infra-adapter 技能）
//
// 设计：每个实现都桥接到本包顶层的 TaskQueue 接口，业务只依赖接口。
package taskqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// asynqQueue 通过 Asynq + Redis 实现 TaskQueue 接口。
type asynqQueue struct {
	client    *asynq.Client
	server    *asynq.Server
	redisAddr string
}

// NewAsynqQueue 构造 Asynq 适配器。redisAddr 格式：`host:port`。
func NewAsynqQueue(redisAddr string) TaskQueue {
	return &asynqQueue{redisAddr: redisAddr}
}

func (q *asynqQueue) ensureClient() *asynq.Client {
	if q.client == nil {
		q.client = asynq.NewClient(asynq.RedisClientOpt{Addr: q.redisAddr})
	}
	return q.client
}

// Enqueue 投递任务到队列。
func (q *asynqQueue) Enqueue(ctx context.Context, topic string, taskID int64, opts EnqueueOptions) error {
	client := q.ensureClient()

	payload, _ := json.Marshal(Payload{TaskID: taskID})
	t := asynq.NewTask(topic, payload)

	var asynqOpts []asynq.Option
	if opts.Queue != "" {
		asynqOpts = append(asynqOpts, asynq.Queue(opts.Queue))
	}
	if opts.MaxRetry > 0 {
		asynqOpts = append(asynqOpts, asynq.MaxRetry(opts.MaxRetry))
	}
	if opts.Delay > 0 {
		asynqOpts = append(asynqOpts, asynq.ProcessIn(opts.Delay))
	}
	if !opts.Deadline.IsZero() {
		asynqOpts = append(asynqOpts, asynq.Deadline(opts.Deadline))
	}
	if opts.Timeout > 0 {
		asynqOpts = append(asynqOpts, asynq.Timeout(opts.Timeout))
	}

	if _, err := client.EnqueueContext(ctx, t, asynqOpts...); err != nil {
		return fmt.Errorf("enqueue topic=%s taskID=%d failed: %w", topic, taskID, err)
	}
	return nil
}

// Consume 启动消费者。
func (q *asynqQueue) Consume(ctx context.Context, subscriptions []Subscription, group string, handler Handler) error {
	// 校验 redisAddr 不为空，空地址会导致运行时连接失败。
	if q.redisAddr == "" {
		return fmt.Errorf("redis address is empty: configure RedisClientOpt.Addr before starting consumer")
	}
	mux := asynq.NewServeMux()
	for _, sub := range subscriptions {
		topic := sub.Pattern
		switch sub.Mode {
		case MatchPrefix:
			topic = sub.Pattern + "*"
		case MatchGlob, MatchRegex:
			// Asynq 不原生支持，需在 handler 内过滤
		}
		mux.HandleFunc(topic, func(ctx context.Context, t *asynq.Task) error {
			var p Payload
			_ = json.Unmarshal(t.Payload(), &p)
			msg := Message{Topic: t.Type(), TaskID: p.TaskID}
			return handler.Handle(ctx, msg)
		})
	}

	q.server = asynq.NewServer(asynq.RedisClientOpt{Addr: q.redisAddr}, asynq.Config{Concurrency: 10})
	return q.server.Start(mux)
}

// Ack 确认消费（asynq handler 返回 nil 即自动 ack，此处为接口一致性保留）。
func (q *asynqQueue) Ack(ctx context.Context, msg Message) error {
	return nil
}

// Nack 否定消费（延迟重试）。
func (q *asynqQueue) Nack(ctx context.Context, msg Message, delay time.Duration) error {
	return fmt.Errorf("nack not directly supported by asynq; use task retry config instead")
}

// Shutdown 优雅停止。
func (q *asynqQueue) Shutdown(ctx context.Context) error {
	if q.server != nil {
		q.server.Shutdown()
	}
	return nil
}
