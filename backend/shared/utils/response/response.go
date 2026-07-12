// Package response 统一 REST 响应结构与写出方式。
//
// 设计（增强错误码支持的统一响应体）：
//   - 统一响应体在此定义，绝不写进 .api 文件（保持 .api 只描述业务 data）。
//   - 与改造版 goctl handler.tpl 配合：handler 成功走 Success(w, resp)，失败走 Error(w, err)。
//   - Error 接受 error：若是 *errorx.CodeError 则带出真实业务码与消息；
//     其余错误统一按内部错误处理，不向客户端泄露实现细节。
//
// 调用链：handler -> Success/Error -> httpx.WriteJson -> 客户端。
package response

import (
	"errors"
	"net/http"

	"GOAI_MODULE/shared/utils/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response 是统一响应结构。
type Response struct {
	Code int         `json:"code"` // 业务码：0 成功；非 0 为 errorx 中定义的错误码。
	Msg  string      `json:"msg"`  // 面向用户的消息。
	Data interface{} `json:"data,omitempty"`
}

// Success 写出成功响应。
func Success(w http.ResponseWriter, data interface{}) {
	httpx.WriteJson(w, http.StatusOK, Response{
		Code: errorx.CodeOK,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 写出成功响应并自定义消息。
func SuccessWithMsg(w http.ResponseWriter, msg string, data interface{}) {
	httpx.WriteJson(w, http.StatusOK, Response{Code: errorx.CodeOK, Msg: msg, Data: data})
}

// Error 写出错误响应。
//   - *errorx.CodeError（含被 %w 包裹的）：保留其 Code 与 Msg。
//   - 其他错误：统一返回内部错误码与安全消息，不泄露实现细节（细节应记日志）。
//
// HTTP 状态统一为 200，客户端依据 body.code 判断业务结果（go-zero 社区通行做法）。
func Error(w http.ResponseWriter, err error) {
	var ce *errorx.CodeError
	if errors.As(err, &ce) {
		httpx.WriteJson(w, http.StatusOK, Response{Code: ce.Code, Msg: ce.Msg})
		return
	}
	httpx.WriteJson(w, http.StatusOK, Response{Code: errorx.CodeInternalError, Msg: "internal error"})
}

// ErrorWithCode 直接以错误码 + 消息写出，供无 error 对象时使用。
func ErrorWithCode(w http.ResponseWriter, code int, msg string) {
	httpx.WriteJson(w, http.StatusOK, Response{Code: code, Msg: msg})
}
