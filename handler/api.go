package handler

import (
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

	// 校验签名
	if wxLogin.Signature != fmt.Sprintf("%x", sha1.Sum([]byte(wxLogin.RawData+wxCode2Session.SessionKey))) {
		return define.ErrorInvalidSignature
	}

	var userID int

	if err := db.MySQL.QueryRow("SELECT UserID FROM user WHERE OpenID = ?", wxCode2Session.OpenID).Scan(&userID); err == sql.ErrNoRows {
		res, err := db.MySQL.Exec("INSERT INTO user (OpenID) VALUES (?)", wxCode2Session.OpenID)
		if err != nil {
			return err
		}

		uid, err := res.LastInsertId()
		if err != nil {
			return err
		}

		userID = int(uid)
	} else if err != nil {
		return err
	}

	redisUserInfo := &define.RedisUserInfo{
		SessionKey: wxCode2Session.SessionKey,
		WxUserInfo: *wxUserInfo,
	}

	if _, err := db.DoOne(db.RedisDefault, "HMSET", redis.Args{}.Add(userID).AddFlat(redisUserInfo)...); err != nil {
		return err
	}

	// 签发令牌
	uidToken.Header["uid"] = userID
	token, err := uidToken.SignedString([]byte(wxCode2Session.SessionKey))
	if err != nil {
		return err
	}

	w.Header().Set("Token", token)

	return nil
}

func serveUpload(w http.ResponseWriter, r *http.Request) error {
	// 校验令牌
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

func servePublish(w http.ResponseWriter, r *http.Request) error {
	// 校验令牌
	token, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	})
	if err != nil {
		return err
	}

	publish := &define.RequestPublish{}
	if err := utils.ReadUnmarshalJSON(r.Body, publish); err != nil {
		return err
	}

	res, err := db.MySQL.Exec("INSERT INTO sku (UserID,Name,Price,MinPrice,Bargain,Intro,Images,WeChatID,Deadline) VALUES (?,?,?,?,?,?,?,?,FROM_UNIXTIME(?))",
		token.Header["uid"], publish.Name, publish.Price, publish.MinPrice, publish.Bargain, publish.Intro, publish.Images, publish.WeChatID, publish.Deadline)
	if err != nil {
		return err
	}

	skuID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Fprintf(w, `{"SkuID":%d}`, skuID)

	return nil
}

func serveShare(w http.ResponseWriter, r *http.Request) error {
	// 校验令牌
	token, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	})
	if err != nil {
		return err
	}

	share := &define.RequestShare{}
	if err := utils.ReadUnmarshalJSON(r.Body, share); err != nil {
		return err
	}

	var shareID int

	if err := db.MySQL.QueryRow("SELECT ShareID FROM share WHERE UserID = ? AND SkuID = ?", token.Header["uid"], share.SkuID).Scan(&shareID); err == sql.ErrNoRows {
		res, err := db.MySQL.Exec("INSERT INTO share (UserID,SkuID) VALUES (?,?)", token.Header["uid"], share.SkuID)
		if err != nil {
			return err
		}

		sid, err := res.LastInsertId()
		if err != nil {
			return err
		}

		shareID = int(sid)
	} else if err != nil {
		return err
	}

	fmt.Fprintf(w, `{"ShareID":%d}`, shareID)

	return nil
}

func serveShow(w http.ResponseWriter, r *http.Request) error {
	show := &define.RequestShow{}
	if err := utils.ReadUnmarshalJSON(r.Body, show); err != nil {
		return err
	}

	rs := &define.ResponseShow{
		Seller: &define.BaseUserInfo{},
	}

	if show.ShareID != 0 {
		rs.Buyer = &define.BaseUserInfo{}

		if err := db.MySQL.QueryRow("SELECT UserID,SkuID FROM share WHERE ShareID = ?", show.ShareID).Scan(&rs.Buyer.UserID, &show.SkuID); err != nil {
			return err
		}

		// 获取买家信息
	}

	if err := db.MySQL.QueryRow("SELECT UserID,Name,Price,ABS(Bargain),Intro,Images,WeChatID,UNIX_TIMESTAMP(Deadline) FROM sku WHERE SkuID = ?", show.SkuID).Scan(&rs.Seller.UserID, &rs.Name, &rs.Price, &rs.Bargain, &rs.Intro, &rs.Images, &rs.WeChatID, &rs.Time); err != nil {
		return err
	}

	// 获取卖家信息

	if rs.Time == 0 {
		rs.Time = -1 // 没有截止时间
	} else if nowUnix := time.Now().Unix(); nowUnix < rs.Time {
		rs.Time -= nowUnix // 倒计时
	} else {
		rs.Time = 0 // 已截止
	}

	if rs.Bargain == 0 {
		rs.Bargain = -1 // 不支持砍价
	} else if show.ShareID != 0 {
		rows, err := db.MySQL.Query("SELECT UserID,BargainPrice FROM bargain WHERE ShareID = ?", show.ShareID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			helper := &define.HelperUserInfo{}

			if err := rows.Scan(&helper.UserID, &helper.BargainPrice); err != nil {
				return err
			}

			// 获取助力者信息

			rs.Helpers = append(rs.Helpers, helper)
		}

		if err := rows.Err(); err != nil {
			return err
		}

		// 剩余砍价机会
		rs.Bargain -= len(rs.Helpers)
	}

	data, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	w.Write(data)

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

	case "/publish":
		err = servePublish(w, r)

	case "/share":
		err = serveShare(w, r)

	case "/show":
		err = serveShow(w, r)

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
