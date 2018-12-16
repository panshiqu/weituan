package define

import (
	"os"

	"github.com/panshiqu/weituan/utils"
)

// GC 全局配置
var GC GlobalConfig

// GlobalConfig 全局配置
type GlobalConfig struct {
	HTTPS     bool
	Address   string // 地址
	KeyFile   string // 私钥
	CertFile  string // 证书
	MySQLDSN  string // MySQL数据源
	RedisAddr string // Redis地址
	RedisAuth string // Redis密码
	AppID     string // 小程序编号
	AppSecret string // 小程序密钥

	StatEverydays  int   // 统计每日
	StatRecentdays []int // 统计最近
}

// Init 初始化
func Init(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return utils.ReadUnmarshalJSON(f, &GC)
}
