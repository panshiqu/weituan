package db

import (
	"database/sql"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/panshiqu/weituan/define"
)

const (
	// RedisDefault 默认数据库
	RedisDefault = 0
)

// MySQL .
var MySQL *sql.DB

// Redis .
var Redis *redis.Pool

// Init 初始化
func Init() (err error) {
	if MySQL, err = sql.Open("mysql", define.GC.MySQLDSN); err != nil {
		return err
	}

	if err = MySQL.Ping(); err != nil {
		return err
	}

	Redis = &redis.Pool{
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

// GetN 获取指定数据库
func GetN(n int) redis.Conn {
	c := Redis.Get()
	c.Do("SELECT", n)
	return c
}

// DoOne 执行一条Redis命令
func DoOne(index int, commandName string, args ...interface{}) (reply interface{}, err error) {
	c := GetN(index)
	defer c.Close()
	return c.Do(commandName, args...)
}
