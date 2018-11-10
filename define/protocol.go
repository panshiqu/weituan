package define

// BaseUserInfo 基础用户信息
type BaseUserInfo struct {
	UserID int `json:",omitempty"` // 用户编号
	WxUserInfo
}

// HelperUserInfo 助力者用户信息
type HelperUserInfo struct {
	BargainInfo
	BaseUserInfo
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

// ShareSkuInfo 分享商品信息
type ShareSkuInfo struct {
	ShareID   int   `json:",omitempty"` // 分享编号
	ShareTime int64 `json:",omitempty"` // 分享时间
	SkuInfo
}

// BargainInfo 砍价信息
type BargainInfo struct {
	BargainTime  int64   `json:",omitempty"` // 时间
	BargainPrice float64 `json:",omitempty"` // 砍价
}

// ResponseUpload 上传
type ResponseUpload struct {
	URL string // 相对地址
}

// RequestPublish 发布
type RequestPublish struct {
	FormID string `json:",omitempty"` // 表单编号（发送模板消息）
	SkuInfo
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
	UserID  int `json:",omitempty"` // 用户编号（客户端分享前不便调用/share接口生成ShareID）
}

// ResponseShow 显示
type ResponseShow struct {
	Seller      *BaseUserInfo     `json:",omitempty"` // 卖家
	Buyer       *BaseUserInfo     `json:",omitempty"` // 买家
	Helpers     []*HelperUserInfo `json:",omitempty"` // 助力者
	CurrentTime int64             `json:",omitempty"` // 当前时间
	SkuInfo
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
	BargainInfo
}

// RequestShareList 分享列表
type RequestShareList struct {
	UserID int `json:",omitempty"` // 用户编号（查看他人分享列表）
}

// ResponseShareList 分享列表
type ResponseShareList struct {
	Buyer *BaseUserInfo   `json:",omitempty"` // 买家
	Skus  []*ShareSkuInfo `json:",omitempty"` // 商品
}
