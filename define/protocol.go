package define

// WxUserInfo 用户信息
type WxUserInfo struct {
	Nickname  string `json:"nickName,omitempty"`  // 昵称
	AvatarURL string `json:"avatarUrl,omitempty"` // 头像
	Gender    int    `json:"gender,omitempty"`    // 性别
}

// WxLogin 登陆
type WxLogin struct {
	Code      string `json:"code,omitempty"`      // 临时登录凭证
	RawData   string `json:"rawData,omitempty"`   // 不包括敏感信息的原始数据字符串，用于计算签名
	Signature string `json:"signature,omitempty"` // 使用 sha1( rawData + sessionkey ) 得到字符串，用于校验用户信息
}

// WxCode2Session 登录凭证校验
type WxCode2Session struct {
	OpenID     string `json:"openid,omitempty"`      // 用户唯一标识
	SessionKey string `json:"session_key,omitempty"` // 会话密钥
	UnionID    string `json:"unionid,omitempty"`     // 用户在开放平台的唯一标识符
	ErrCode    int    `json:"errcode,omitempty"`     // 错误码
	ErrMsg     string `json:"errMsg,omitempty"`      // 错误信息
}