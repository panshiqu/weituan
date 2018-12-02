package define

import "fmt"

const (
	// ErrSuccess 成功
	ErrSuccess int = 0

	// ErrFailure 失败
	ErrFailure int = 1

	// ErrUnsupportedAPI 不支持的接口
	ErrUnsupportedAPI int = 2

	// ErrInvalidSignature 无效的签名
	ErrInvalidSignature int = 3

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken int = 4
)

var (
	// ErrorSuccess 成功
	ErrorSuccess = &MyError{ErrCode: ErrSuccess, ErrDesc: "success"}

	// ErrorFailure 失败
	ErrorFailure = &MyError{ErrCode: ErrFailure, ErrDesc: "failure"}

	// ErrorUnsupportedAPI 不支持的接口
	ErrorUnsupportedAPI = &MyError{ErrCode: ErrUnsupportedAPI, ErrDesc: "unsupported api"}

	// ErrorInvalidSignature 无效的签名
	ErrorInvalidSignature = &MyError{ErrCode: ErrInvalidSignature, ErrDesc: "invalid signature"}

	// ErrorInvalidToken 无效的令牌
	ErrorInvalidToken = &MyError{ErrCode: ErrInvalidToken, ErrDesc: "invalid token"}
)

// MyError 错误
type MyError struct {
	ErrCode int    `json:",omitempty"` // 错误码
	ErrDesc string `json:",omitempty"` // 错误描述
}

func (m *MyError) Error() string {
	return fmt.Sprintf(`{"ErrCode":%d,"ErrDesc":"%s"}`, m.ErrCode, m.ErrDesc)
}

// NewFailure 失败
func NewFailure(desc string) *MyError {
	return &MyError{ErrCode: ErrFailure, ErrDesc: desc}
}

// NewError 错误
func NewError(code int, desc string) *MyError {
	return &MyError{ErrCode: code, ErrDesc: desc}
}

// WxError 微信错误
type WxError struct {
	ErrCode int    `json:"errcode,omitempty"` // 错误码
	ErrMsg  string `json:"errMsg,omitempty"`  // 错误信息
}
