// Package jwtx 提供 JWT 的签发、刷新与解析，以及登录身份 Identity。
//
// 分工：go-zero REST 内建 jwt 中间件负责校验签名与过期；本包用于
//   - 业务侧签发 access / refresh token（登录 logic）；
//   - 解析 token 得到 Identity（中间件写入上下文，见 shared/utils/contexts）。
//
// 调用链：登录 logic -> GetJwtToken -> 返回 token；
//         请求 -> jwt 中间件校验 -> ParseToken -> contexts.SetIdentity -> logic 读取。
package jwtx

import "github.com/golang-jwt/jwt/v4"

// Config 收口 JWT 配置。
type Config struct {
	AccessSecret string `json:",optional"`       // AccessSecret 来自安全配置，禁止硬编码、入日志。
	AccessExpire int64  `json:",default=86400"`  // AccessExpire 是 access token 过期秒数。
	RefreshAfter int64  `json:",default=604800"` // RefreshAfter 是 refresh token 过期秒数。
}

type (
	// UserClaims 是携带自定义身份的 JWT claims。
	UserClaims struct {
		jwt.RegisteredClaims
		CustomClaims Identity `json:"identity"`
	}

	// Identity 是登录用户身份，随 token 下发、经中间件写入请求上下文。
	Identity struct {
		UserID          int64  `json:"user_id"`
		DeviceSessionID int64  `json:"device_session_id"`
		RoleCode        string `json:"role_code"`
	}
)

// GetJwtToken 生成携带 Identity 的 access token。
func GetJwtToken(secretKey string, iat, seconds int64, claims Identity) (string, error) {
	mapClaims := make(jwt.MapClaims)
	mapClaims["exp"] = iat + seconds
	mapClaims["iat"] = iat
	mapClaims["identity"] = claims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(secretKey))
}

// GetRefreshToken 生成 refresh token（当前复用签发逻辑，可按需换密钥或过期）。
func GetRefreshToken(secretKey string, iat, seconds int64, claims Identity) (string, error) {
	return GetJwtToken(secretKey, iat, seconds, claims)
}

// ParseToken 校验并解析 token，返回 Identity。
func ParseToken(tokenString string, secretKey string) (*Identity, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return &claims.CustomClaims, nil
	}
	return nil, jwt.ErrInvalidKey
}
