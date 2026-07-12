package contexts

import (
	"context"

	"GOAI_MODULE/shared/pkg/jwtx"
)

type contextKey string

const (
	identityKey  contextKey = "identity"
	requestIDKey contextKey = "request_id"
)

// SetIdentity 将认证后的用户身份写入请求上下文。
// 调用链：UserAuthMiddleware -> contexts.SetIdentity -> account logic。
func SetIdentity(ctx context.Context, identity *jwtx.Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

// GetIdentity 从请求上下文读取用户身份。
// 仅认证中间件保护过的接口可直接依赖该方法；未认证上下文会返回 nil。
func GetIdentity(ctx context.Context) *jwtx.Identity {
	identity, _ := ctx.Value(identityKey).(*jwtx.Identity)
	return identity
}

// GetUserID 从请求上下文读取当前登录用户 ID。
// logic 层必须通过该方法读取用户 ID，避免散落 context key 字符串。
func GetUserID(ctx context.Context) int64 {
	identity := GetIdentity(ctx)
	if identity == nil {
		return 0
	}
	return identity.UserID
}

// GetRequestID 从请求上下文读取 request id。
func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}

// SetRequestID 将 request id 写入请求上下文。
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}
