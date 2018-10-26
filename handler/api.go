package handler

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/panshiqu/weituan/define"
	"github.com/panshiqu/weituan/utils"
)

func serveLogin(w http.ResponseWriter, r *http.Request) error {
	wxLogin := &define.WxLogin{}
	if err := utils.ReadUnmarshalJSON(r.Body, wxLogin); err != nil {
		return err
	}

	wxUserInfo := &define.WxUserInfo{}
	if err := json.Unmarshal([]byte(wxLogin.RawData), wxUserInfo); err != nil {
		return err
	}

	// 登录凭证校验
	wxCode2Session := &define.WxCode2Session{}
	if err := utils.HTTPGetJSON(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", define.GC.AppID, define.GC.AppSecret, wxLogin.Code), wxCode2Session); err != nil {
		return err
	}

	if wxCode2Session.ErrCode != 0 {
		return define.NewFailure(fmt.Sprintf("%d:%s", wxCode2Session.ErrCode, wxCode2Session.ErrMsg))
	}

	// 计算比对签名
	if wxLogin.Signature != fmt.Sprintf("%x", sha1.Sum([]byte(wxLogin.RawData+wxCode2Session.SessionKey))) {
		return define.ErrorInvalidSignature
	}

	return nil
}

// ServeHTTP .
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.URL.Path {
	case "/login":
		err = serveLogin(w, r)

	default:
		err = define.ErrorUnsupportedAPI
	}

	if _, ok := err.(*define.MyError); ok {
		fmt.Fprint(w, err)
	}

	log.Println(r.URL.Path, err)
}

// ServeFiles .
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	name := "." + r.URL.Path

	fi, err := os.Stat(name)
	if err != nil {
		log.Println("ServeFile", err)
		http.NotFound(w, r)
		return
	}

	if !fi.Mode().IsRegular() {
		log.Println("ServeFile not regular", name)
		http.NotFound(w, r)
		return
	}

	f, err := os.Open(name)
	if err != nil {
		log.Println("ServeFile", err)
		http.NotFound(w, r)
		return
	}

	defer f.Close()

	io.Copy(w, f)
}
