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
