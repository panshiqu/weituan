package define

const (
	// PublishTemplateID 发布模板编号
	PublishTemplateID = "m5-VrAq0-h1nTEMXGiC1yGV9cyzAHB1gJtWitCvHklE"
)

// WxUserInfo 用户信息
type WxUserInfo struct {
	Nickname  string `json:"nickName,omitempty" redis:"nickName"`   // 昵称
	AvatarURL string `json:"avatarUrl,omitempty" redis:"avatarUrl"` // 头像
	Gender    int    `json:"gender,omitempty" redis:"gender"`       // 性别
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
	WxError
}

// WxAccessToken 访问凭证
type WxAccessToken struct {
	AccessToken string `json:"access_token,omitempty"` // 获取到的凭证
	ExpiresIn   int    `json:"expires_in,omitempty"`   // 凭证有效时间，单位：秒
	WxError
}

// WxTemplateMessage 模板消息
type WxTemplateMessage struct {
	ToUser          string                      `json:"touser,omitempty"`           // 接收者（用户）的 openid
	TemplateID      string                      `json:"template_id,omitempty"`      // 所需下发的模板消息的id
	Page            string                      `json:"page,omitempty"`             // 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	FormID          string                      `json:"form_id,omitempty"`          // 表单提交场景下，为 submit 事件带上的 formId；支付场景下，为本次支付的 prepay_id
	Data            map[string]*WxTemplateValue `json:"data,omitempty"`             // 模板内容，不填则下发空模板
	EmphasisKeyword string                      `json:"emphasis_keyword,omitempty"` // 模板需要放大的关键词，不填则默认无放大
}

// WxTemplateValue 模板值
type WxTemplateValue struct {
	Value string `json:"value,omitempty"`
}
