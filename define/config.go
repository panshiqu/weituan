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
	MysqlDSN  string // MYSQL数据源
	RedisAddr string // REDIS地址
	RedisAuth string // REDIS密码
	AppID     string // 小程序编号
	AppSecret string // 小程序密钥
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
