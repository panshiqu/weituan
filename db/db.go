package db

import (
	"database/sql"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/panshiqu/weituan/define"
)

// MYSQL .
var MYSQL *sql.DB

// REDIS .
var REDIS *redis.Pool

// Init 初始化
func Init() (err error) {
	if MYSQL, err = sql.Open("mysql", define.GC.MysqlDSN); err != nil {
		return err
	}

	if err := MYSQL.Ping(); err != nil {
		return err
	}

	REDIS = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", define.GC.RedisAddr)
			if err != nil {
				return nil, err
			}
			if define.GC.RedisAuth != "" {
				if _, err := c.Do("AUTH", define.GC.RedisAuth); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}
