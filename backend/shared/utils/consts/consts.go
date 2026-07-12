// Package consts 汇总跨域复用的常量：状态字符串、上下文 key、通用枚举。
//
// 约定：状态用显式字符串（pending/running/success/failed），不用魔法数字；
// 域内私有的常量放各自 model/logic 包，不堆到这里。
package consts

// 通用任务 / 资源状态（显式字符串，前后端一致）。
const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

// 请求上下文 key（由中间件写入，logic 读取）。
type ctxKey string

const (
	CtxUserID ctxKey = "uid" // CtxUserID 是 jwt 中间件解析出的用户 ID。
)
