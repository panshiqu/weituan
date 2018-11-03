package define

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
	ErrCode    int    `json:"errcode,omitempty"`     // 错误码
	ErrMsg     string `json:"errMsg,omitempty"`      // 错误信息
}

// RedisUserInfo 缓存用户信息
type RedisUserInfo struct {
	SessionKey string
	WxUserInfo
}

// ResponseUpload 上传
type ResponseUpload struct {
	URL string // 相对地址
}

// RequestPublish 发布
type RequestPublish struct {
	FormID   string  `json:",omitempty"` // 表单编号（发送模板消息）
	SkuID    int     `json:",omitempty"` // 商品编号（修改时提供）
	Name     string  `json:",omitempty"` // 名称
	Price    float64 `json:",omitempty"` // 价格
	MinPrice float64 `json:",omitempty"` // 底价
	Bargain  int     `json:",omitempty"` // 砍价（0不支持砍价 +n随机砍N次 -n等值砍N次）
	Intro    string  `json:",omitempty"` // 介绍
	Images   string  `json:",omitempty"` // 图片
	WeChatID string  `json:",omitempty"` // 微信号（卖家）
	Deadline int64   `json:",omitempty"` // 截止时间
}

// ResponsePublish 发布
type ResponsePublish struct {
	SkuID int `json:",omitempty"` // 商品编号
}

// RequestShare 分享
type RequestShare struct {
	SkuID int `json:",omitempty"` // 商品编号
}

// ResponseShare 分享
type ResponseShare struct {
	ShareID int `json:",omitempty"` // 分享编号
}

// RequestShow 显示
type RequestShow struct {
	SkuID   int `json:",omitempty"` // 商品编号
	ShareID int `json:",omitempty"` // 分享编号
}

// BaseUserInfo 基础用户信息
type BaseUserInfo struct {
	UserID int `json:",omitempty"` // 用户编号
	WxUserInfo
}

// HelperUserInfo 助力者用户信息
type HelperUserInfo struct {
	BargainPrice float64 `json:",omitempty"` // 砍价
	BaseUserInfo
}

// ResponseShow 显示
type ResponseShow struct {
	Seller      *BaseUserInfo     `json:",omitempty"` // 卖家
	Buyer       *BaseUserInfo     `json:",omitempty"` // 买家
	Helpers     []*HelperUserInfo `json:",omitempty"` // 助力者
	CurrentTime int64             `json:",omitempty"` // 当前时间
	SkuInfo
}

// SkuInfo 商品信息
type SkuInfo struct {
	SkuID       int     `json:",omitempty"` // 商品编号
	Name        string  `json:",omitempty"` // 名称
	Price       float64 `json:",omitempty"` // 价格
	MinPrice    float64 `json:",omitempty"` // 底价
	Bargain     int     `json:",omitempty"` // 砍价（0不支持砍价 +n随机砍N次 -n等值砍N次）
	Intro       string  `json:",omitempty"` // 介绍
	Images      string  `json:",omitempty"` // 图片
	WeChatID    string  `json:",omitempty"` // 微信号（卖家）
	Deadline    int64   `json:",omitempty"` // 截止时间
	PublishTime int64   `json:",omitempty"` // 发布时间
	Status      int     `json:",omitempty"` // 状态（审核）（暂未实现）
}

// RequestList 列表
type RequestList struct {
	UserID int `json:",omitempty"` // 用户编号（查看他人列表）
}

// ResponseList 列表
type ResponseList struct {
	Seller *BaseUserInfo `json:",omitempty"` // 卖家
	Skus   []*SkuInfo    `json:",omitempty"` // 商品
}

// RequestBargain 砍价
type RequestBargain struct {
	ShareID int `json:",omitempty"` // 分享编号
}

// ResponseBargain 砍价
type ResponseBargain struct {
	BargainPrice float64 `json:",omitempty"` // 砍价
}
