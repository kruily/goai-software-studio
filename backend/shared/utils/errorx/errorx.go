// Package errorx 定义业务错误类型与集中式错误码表。
//
// 设计（结合 go-zero 官方 CodeError 与社区 looklook 的错误码表模式）：
//   - CodeError 携带业务码 + 用户可见消息，实现 error 接口。
//   - 错误码集中登记在此，配 message 表，避免散落各处、便于前端对照。
//   - logic 只返回 errorx.New(code, msg) 或 errorx.NewByCode(code)；
//     response.Error 会用 errors.As 提取（支持 %w 包裹）。
//
// 约定：不把内部实现细节、堆栈、SQL 泄露给用户；排查细节写日志，返回给用户的用 Msg。
package errorx

import "fmt"

// CodeError 是带错误码的业务错误。
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error 实现 error 接口。
func (e *CodeError) Error() string {
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Msg)
}

// New 构造带自定义消息的业务错误。
func New(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

// NewByCode 用错误码构造错误，消息取自集中式错误码表。
func NewByCode(code int) *CodeError {
	return &CodeError{Code: code, Msg: MapMsg(code)}
}

// 错误码约定（各项目按域扩展；保持语义稳定，前后端共享）。
const (
	CodeOK            = 0
	CodeInvalidParam  = 40000 // 请求参数非法。
	CodeUnauthorized  = 40100 // 未认证或 token 失效。
	CodeForbidden     = 40300 // 无权限。
	CodeNotFound      = 40400 // 资源不存在。
	CodeConflict      = 40900 // 资源冲突（如唯一键重复）。
	CodeTooManyReq    = 42900 // 请求过于频繁。
	CodeInternalError = 50000 // 服务内部错误。
)

// messages 是错误码到默认消息的集中映射。业务新增码时在此登记。
var messages = map[int]string{
	CodeOK:            "success",
	CodeInvalidParam:  "invalid parameter",
	CodeUnauthorized:  "unauthorized",
	CodeForbidden:     "forbidden",
	CodeNotFound:      "not found",
	CodeConflict:      "conflict",
	CodeTooManyReq:    "too many requests",
	CodeInternalError: "internal error",
}

// MapMsg 返回错误码对应的默认消息，未登记时回退为内部错误消息。
func MapMsg(code int) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return messages[CodeInternalError]
}
