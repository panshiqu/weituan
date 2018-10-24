package define

import "fmt"

const (
	// ErrSuccess 成功
	ErrSuccess int = 0

	// ErrFailure 失败
	ErrFailure int = 1

	// ErrUnsupportedAPI 不支持的接口
	ErrUnsupportedAPI int = 2
)

var (
	// ErrorSuccess 成功
	ErrorSuccess = &MyError{Errcode: ErrSuccess, Errdesc: "success"}

	// ErrorFailure 失败
	ErrorFailure = &MyError{Errcode: ErrFailure, Errdesc: "failure"}

	// ErrorUnsupportedAPI 不支持的接口
	ErrorUnsupportedAPI = &MyError{Errcode: ErrUnsupportedAPI, Errdesc: "unsupported api"}
)

// MyError 错误
type MyError struct {
	Errcode int    `json:",omitempty"` // 错误码
	Errdesc string `json:",omitempty"` // 错误描述
}

func (m *MyError) Error() string {
	return fmt.Sprintf(`{"Errcode":%d,"Errdesc":"%s"}`, m.Errcode, m.Errdesc)
}