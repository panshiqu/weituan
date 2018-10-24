package define

import (
	"os"

	"github.com/panshiqu/weituan/utils"
)

// GC 全局配置
var GC GlobalConfig

// GlobalConfig 全局配置
type GlobalConfig struct {
	HTTPS    bool
	Address  string // 地址
	KeyFile  string // 私钥
	CertFile string // 证书
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
