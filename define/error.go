package define

import "fmt"

const (
	// ErrcodeSuccess 成功
	ErrcodeSuccess int = 0

	// ErrcodeFailure 失败
	ErrcodeFailure int = 1

	// ErrcodeUnsupportedAPI 不支持的接口
	ErrcodeUnsupportedAPI int = 2
)

var (
	// ErrSuccess 成功
	ErrSuccess = &MyError{Errcode: ErrcodeSuccess, Errdesc: "success"}

	// ErrFailure 失败
	ErrFailure = &MyError{Errcode: ErrcodeFailure, Errdesc: "failure"}

	// ErrUnsupportedAPI 不支持的接口
	ErrUnsupportedAPI = &MyError{Errcode: ErrcodeUnsupportedAPI, Errdesc: "unsupported api"}
)

// MyError 错误
type MyError struct {
	Errcode int    `json:",omitempty"` // 错误码
	Errdesc string `json:",omitempty"` // 错误描述
}

func (m *MyError) Error() string {
	return fmt.Sprintf(`{"Errcode":%d,"Errdesc":"%s"}`, m.Errcode, m.Errdesc)
}
