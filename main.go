package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/panshiqu/weituan/db"
	"github.com/panshiqu/weituan/define"
	"github.com/panshiqu/weituan/handler"
)

var conf = flag.String("conf", "./conf.json", "conf")

func main() {
	flag.Parse()

	// 配置初始化
	if err := define.Init(*conf); err != nil {
		log.Fatal(err)
	}

	// 数据库初始化
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler.ServeHTTP)
	http.HandleFunc("/files/", handler.ServeFiles)

	if !define.GC.HTTPS {
		log.Fatal(http.ListenAndServe(define.GC.Address, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(define.GC.Address, define.GC.CertFile, define.GC.KeyFile, nil))
	}
}
