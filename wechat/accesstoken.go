package wechat

import (
	"fmt"
	"sync"
	"time"

	"github.com/panshiqu/weituan/define"
	"github.com/panshiqu/weituan/utils"
)

var accessToken string
var accessTokenMutex sync.Mutex

// AccessToken 获取访问凭证
func AccessToken() (string, error) {
	accessTokenMutex.Lock()
	defer accessTokenMutex.Unlock()

	if accessToken != "" {
		return accessToken, nil
	}

	wxAccessToken := &define.WxAccessToken{}
	if err := utils.HTTPGetJSON(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", define.GC.AppID, define.GC.AppSecret), wxAccessToken); err != nil {
		return "", err
	}

	if wxAccessToken.ErrCode != 0 {
		return "", define.NewFailure(fmt.Sprintf("%d:%s", wxAccessToken.ErrCode, wxAccessToken.ErrMsg))
	}

	time.AfterFunc(time.Duration(wxAccessToken.ExpiresIn)*time.Second, func() {
		accessToken = ""
	})

	accessToken = wxAccessToken.AccessToken

	return accessToken, nil
}
