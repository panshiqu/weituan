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
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/panshiqu/weituan/db"
	"github.com/panshiqu/weituan/define"
	"github.com/panshiqu/weituan/utils"
)

// uidToken 令牌
var uidToken = jwt.New(jwt.SigningMethodHS256)

func updateUserInfo(uid interface{}, key, data, signature string) error {
	wxUserInfo := &define.WxUserInfo{}
	if err := json.Unmarshal([]byte(data), wxUserInfo); err != nil {
		return err
	}

	// 校验签名
	if signature != fmt.Sprintf("%x", sha1.Sum([]byte(data+key))) {
		return define.ErrorInvalidSignature
	}

	if _, err := db.MySQL.Exec("UPDATE user SET Nickname=?, AvatarURL=?, Gender=? WHERE UserID = ?", wxUserInfo.Nickname, wxUserInfo.AvatarURL, wxUserInfo.Gender, uid); err != nil {
		return err
	}

	if _, err := db.DoOne(db.RedisDefault, "HMSET", redis.Args{}.Add(uid).AddFlat(wxUserInfo)...); err != nil {
		return err
	}

	return nil
}

func serveLogin(w http.ResponseWriter, r *http.Request) error {
	wxLogin := &define.WxLogin{}
	if err := utils.ReadUnmarshalJSON(r.Body, wxLogin); err != nil {
		return err
	}

	// 客户端要求拆分接口
	if wxLogin.Code == "" {
		// 校验令牌
		token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
			return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
		})
		if err != nil {
			return define.ErrorInvalidToken
		}

		key, err := redis.String(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
		if err != nil {
			return err
		}

		return updateUserInfo(token.Header["uid"], key, wxLogin.RawData, wxLogin.Signature)
	}

	// 登录凭证校验
	wxCode2Session := &define.WxCode2Session{}
	if err := utils.HTTPGetJSON(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", define.GC.AppID, define.GC.AppSecret, wxLogin.Code), wxCode2Session); err != nil {
		return err
	}

	if wxCode2Session.ErrCode != 0 {
		return define.NewFailure(fmt.Sprintf("%d:%s", wxCode2Session.ErrCode, wxCode2Session.ErrMsg))
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

	// 客户端要求拆分接口
	if wxLogin.Signature != "" {
		if err := updateUserInfo(userID, wxCode2Session.SessionKey, wxLogin.RawData, wxLogin.Signature); err != nil {
			return err
		}
	}

	if _, err := db.DoOne(db.RedisDefault, "HSET", userID, "SessionKey", wxCode2Session.SessionKey); err != nil {
		return err
	}

	// 签发令牌
	uidToken.Header["uid"] = userID
	token, err := uidToken.SignedString([]byte(wxCode2Session.SessionKey))
	if err != nil {
		return err
	}

	w.Header().Set("token", token)

	return nil
}

func serveUpload(w http.ResponseWriter, r *http.Request) error {
	// 校验令牌
	if _, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	}); err != nil {
		return define.ErrorInvalidToken
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
	token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	})
	if err != nil {
		return define.ErrorInvalidToken
	}

	publish := &define.RequestPublish{}
	if err := utils.ReadUnmarshalJSON(r.Body, publish); err != nil {
		return err
	}

	if publish.SkuID == 0 {
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
	} else {
		if _, err := db.MySQL.Exec("UPDATE sku SET Name=?, Price=?, MinPrice=?, Bargain=?, Intro=?, Images=?, WeChatID=?, Deadline=FROM_UNIXTIME(?) WHERE SkuID = ? AND UserID = ?",
			publish.Name, publish.Price, publish.MinPrice, publish.Bargain, publish.Intro, publish.Images, publish.WeChatID, publish.Deadline, publish.SkuID, token.Header["uid"]); err != nil {
			return err
		}
	}

	return nil
}

func serveShare(w http.ResponseWriter, r *http.Request) error {
	// 校验令牌
	token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	})
	if err != nil {
		return define.ErrorInvalidToken
	}

	share := &define.RequestShare{}
	if err := utils.ReadUnmarshalJSON(r.Body, share); err != nil {
		return err
	}

	id, err := doShare(token.Header["uid"], share.SkuID)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, `{"ShareID":%d}`, id)

	return nil
}

func serveShow(w http.ResponseWriter, r *http.Request) error {
	show := &define.RequestShow{}
	if err := utils.ReadUnmarshalJSON(r.Body, show); err != nil {
		return err
	}

	rs := &define.ResponseShow{
		Seller:      &define.BaseUserInfo{},
		CurrentTime: time.Now().Unix(),
	}

	if show.UserID != 0 {
		id, err := doShare(show.UserID, show.SkuID)
		if err != nil {
			return err
		}

		rs.Buyer = &define.BaseUserInfo{
			UserID: show.UserID,
		}

		show.ShareID = id
	} else if show.ShareID != 0 {
		rs.Buyer = &define.BaseUserInfo{}

		if err := db.MySQL.QueryRow("SELECT UserID,SkuID FROM share WHERE ShareID = ?", show.ShareID).Scan(&rs.Buyer.UserID, &show.SkuID); err != nil {
			return err
		}
	}

	if rs.Buyer != nil {
		// 获取买家信息
		if err := getWxUserInfo(rs.Buyer); err != nil {
			return err
		}
	}

	if err := db.MySQL.QueryRow("SELECT SkuID,UserID,Name,Price,MinPrice,Bargain,Intro,Images,WeChatID,UNIX_TIMESTAMP(Deadline),UNIX_TIMESTAMP(PublishTime) FROM sku WHERE SkuID = ?", show.SkuID).Scan(&rs.SkuID, &rs.Seller.UserID, &rs.Name, &rs.Price, &rs.MinPrice, &rs.Bargain, &rs.Intro, &rs.Images, &rs.WeChatID, &rs.Deadline, &rs.PublishTime); err != nil {
		return err
	}

	// 获取卖家信息
	if err := getWxUserInfo(rs.Seller); err != nil {
		return err
	}

	if rs.Bargain != 0 && show.ShareID != 0 {
		rows, err := db.MySQL.Query("SELECT UserID,BargainPrice,UNIX_TIMESTAMP(BargainTime) FROM bargain WHERE ShareID = ?", show.ShareID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			helper := &define.HelperUserInfo{}

			if err := rows.Scan(&helper.UserID, &helper.BargainPrice, &helper.BargainTime); err != nil {
				return err
			}

			// 获取助力者信息
			if err := getWxUserInfo(&helper.BaseUserInfo); err != nil {
				return err
			}

			rs.Helpers = append(rs.Helpers, helper)
		}

		if err := rows.Err(); err != nil {
			return err
		}
	}

	data, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	w.Write(data)

	return nil
}

func serveList(w http.ResponseWriter, r *http.Request) error {
	seller := &define.BaseUserInfo{}

	if r.ContentLength == 0 {
		// 校验令牌
		token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
			return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
		})
		if err != nil {
			return define.ErrorInvalidToken
		}

		uid, err := strconv.Atoi(fmt.Sprint(token.Header["uid"]))
		if err != nil {
			return err
		}

		seller.UserID = uid
	} else {
		list := &define.RequestList{}
		if err := utils.ReadUnmarshalJSON(r.Body, list); err != nil {
			return err
		}

		seller.UserID = list.UserID
	}

	// 获取卖家信息
	if err := getWxUserInfo(seller); err != nil {
		return err
	}

	rl := &define.ResponseList{
		Seller: seller,
	}

	rows, err := db.MySQL.Query("SELECT SkuID,Name,Price,MinPrice,Bargain,Intro,Images,WeChatID,UNIX_TIMESTAMP(Deadline),UNIX_TIMESTAMP(PublishTime) FROM sku WHERE UserID = ?", rl.Seller.UserID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		sku := &define.SkuInfo{}

		if err := rows.Scan(&sku.SkuID, &sku.Name, &sku.Price, &sku.MinPrice, &sku.Bargain, &sku.Intro, &sku.Images, &sku.WeChatID, &sku.Deadline, &sku.PublishTime); err != nil {
			return err
		}

		rl.Skus = append(rl.Skus, sku)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	data, err := json.Marshal(rl)
	if err != nil {
		return err
	}

	w.Write(data)

	return nil
}

func serveBargain(w http.ResponseWriter, r *http.Request) error {
	// 校验令牌
	token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
		return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
	})
	if err != nil {
		return define.ErrorInvalidToken
	}

	bargain := &define.RequestBargain{}
	if err := utils.ReadUnmarshalJSON(r.Body, bargain); err != nil {
		return err
	}

	var skuID int

	if err := db.MySQL.QueryRow("SELECT SkuID FROM share WHERE ShareID = ?", bargain.ShareID).Scan(&skuID); err != nil {
		return err
	}

	var m float64 // 可砍总额
	var n int     // 可砍次数

	if err := db.MySQL.QueryRow("SELECT Price-MinPrice,Bargain FROM sku WHERE SkuID = ?", skuID).Scan(&m, &n); err != nil {
		return err
	}

	// 不支持砍价
	if n == 0 {
		return &define.MyError{100001, "不支持砍价"}
	}

	var a float64 // 已砍额度
	var b int     // 已砍次数

	if err := db.MySQL.QueryRow("SELECT IFNULL(SUM(BargainPrice),0),COUNT(UserID) FROM bargain WHERE ShareID = ?", bargain.ShareID).Scan(&a, &b); err != nil {
		return err
	}

	// 已砍到底价
	if a >= m {
		return &define.MyError{100002, "已砍到底价"}
	}

	// 已砍够次数
	if b >= utils.AbsInt(n) {
		return &define.MyError{100003, "已砍够次数"}
	}

	var v float64

	// 等值砍
	if n < 0 {
		v = (m - a) / float64(-n-b)
	} else {
		return define.NewFailure("暂不支持随机砍")
	}

	if _, err := db.MySQL.Exec("INSERT INTO bargain (ShareID,UserID,BargainPrice) VALUES (?,?,?)", bargain.ShareID, token.Header["uid"], v); err != nil {
		return err
	}

	fmt.Fprintf(w, `{"BargainTime":%d,"BargainPrice":%.2f}`, time.Now().Unix(), v)

	return nil
}

func doShare(userID interface{}, skuID int) (int, error) {
	var shareID int

	if err := db.MySQL.QueryRow("SELECT ShareID FROM share WHERE UserID = ? AND SkuID = ?", userID, skuID).Scan(&shareID); err == sql.ErrNoRows {
		res, err := db.MySQL.Exec("INSERT INTO share (UserID,SkuID) VALUES (?,?)", userID, skuID)
		if err != nil {
			return 0, err
		}

		sid, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}

		return int(sid), nil
	} else if err != nil {
		return 0, err
	}

	return shareID, nil
}

func getWxUserInfo(info *define.BaseUserInfo) error {
	v, err := redis.Values(db.DoOne(db.RedisDefault, "HGETALL", info.UserID))
	if err != nil {
		return err
	}

	if len(v) != 0 {
		return redis.ScanStruct(v, info)
	}

	return db.MySQL.QueryRow("SELECT Nickname,AvatarURL,Gender FROM user WHERE UserID = ?", info.UserID).Scan(&info.Nickname, &info.AvatarURL, &info.Gender)
}

func serveShareList(w http.ResponseWriter, r *http.Request) error {
	buyer := &define.BaseUserInfo{}

	if r.ContentLength == 0 {
		// 校验令牌
		token, err := jwt.Parse(r.Header.Get("token"), func(token *jwt.Token) (interface{}, error) {
			return redis.Bytes(db.DoOne(db.RedisDefault, "HGET", token.Header["uid"], "SessionKey"))
		})
		if err != nil {
			return define.ErrorInvalidToken
		}

		uid, err := strconv.Atoi(fmt.Sprint(token.Header["uid"]))
		if err != nil {
			return err
		}

		buyer.UserID = uid
	} else {
		shareList := &define.RequestShareList{}
		if err := utils.ReadUnmarshalJSON(r.Body, shareList); err != nil {
			return err
		}

		buyer.UserID = shareList.UserID
	}

	// 获取买家信息
	if err := getWxUserInfo(buyer); err != nil {
		return err
	}

	rsl := &define.ResponseShareList{
		Buyer: buyer,
	}

	rows, err := db.MySQL.Query("SELECT ShareID,SkuID,UNIX_TIMESTAMP(ShareTime) FROM share WHERE UserID = ?", rsl.Buyer.UserID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		sku := &define.ShareSkuInfo{}

		if err := rows.Scan(&sku.ShareID, &sku.SkuID, &sku.ShareTime); err != nil {
			return err
		}

		if err := db.MySQL.QueryRow("SELECT Name,Price,MinPrice,Bargain,Intro,Images,WeChatID,UNIX_TIMESTAMP(Deadline),UNIX_TIMESTAMP(PublishTime) FROM sku WHERE SkuID = ?", sku.SkuID).Scan(&sku.Name, &sku.Price, &sku.MinPrice, &sku.Bargain, &sku.Intro, &sku.Images, &sku.WeChatID, &sku.Deadline, &sku.PublishTime); err != nil {
			return err
		}

		rsl.Skus = append(rsl.Skus, sku)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	data, err := json.Marshal(rsl)
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

	case "/list":
		err = serveList(w, r)

	case "/bargain":
		err = serveBargain(w, r)

	case "/shareList":
		err = serveShareList(w, r)

	default:
		err = define.ErrorUnsupportedAPI
	}

	if _, ok := err.(*define.MyError); ok {
		fmt.Fprint(w, err)
	} else if err != nil {
		fmt.Fprint(w, define.NewFailure(err.Error()))
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
