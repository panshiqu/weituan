package db

import (
	"database/sql"

	"github.com/panshiqu/weituan/define"
)

// MYSQL .
var MYSQL *sql.DB

// Init 初始化
func Init() (err error) {
	if MYSQL, err = sql.Open("mysql", define.GC.MysqlDSN); err != nil {
		return err
	}

	if err := MYSQL.Ping(); err != nil {
		return err
	}

	return nil
}
