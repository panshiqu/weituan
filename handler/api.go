package handler

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/panshiqu/weituan/db"
	"github.com/panshiqu/weituan/define"
	"github.com/panshiqu/weituan/utils"
)

// uidToken 令牌
var uidToken = jwt.New(jwt.SigningMethodHS256)

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

	// 签发令牌
	uidToken.Header["uid"] = 123
	token, err := uidToken.SignedString([]byte("hello"))
	if err != nil {
		return err
	}

	w.Header().Set("Token", token)

	return nil
}

func serveUpload(w http.ResponseWriter, r *http.Request) error {
	if _, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	}); err != nil {
		return err
	}

	f, fh, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	// md5命名
	fn := fmt.Sprintf("./files/%x%s", h.Sum(nil), filepath.Ext(fh.Filename))

	fmt.Fprintf(w, `{"URL":"%s"}`, fn[1:])

	// 文件已存在
	if _, err := os.Stat(fn); err == nil {
		return nil
	}

	n, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer n.Close()

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := io.Copy(n, f); err != nil {
		return err
	}

	return nil
}

// ServeHTTP .
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.URL.Path {
	case "/login":
		err = serveLogin(w, r)

	case "/upload":
		err = serveUpload(w, r)

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
