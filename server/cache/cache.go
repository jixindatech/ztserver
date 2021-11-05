package cache

import (
	"fmt"
	redis "github.com/gomodule/redigo/redis"
	"time"
	"zt-server/settings"
)

var redisConfig *settings.Redis
var pool *redis.Pool

func Init(redisSetting *settings.Redis) error {
	redisConfig = redisSetting

	pool = &redis.Pool{
		MaxIdle:     10, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   20, //最大数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConfig.Addr)
			if err != nil {
				return nil, err
			}
			if redisConfig.Password != "" {
				if _, err := c.Do("AUTH", redisConfig.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if redisConfig.Db != 0 {
				if _, err := c.Do("SELECT", redisConfig.Db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	client := pool.Get()
	_, err := client.Do("PING")
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	client.Close()
	return nil
}

func Set(key, value string) error {
	client := pool.Get()
	_, err := client.Do("SET", key, value)
	if err != nil {
		return err
	}

	client.Close()
	return nil
}

func Get(key string) (interface{}, error) {
	client := pool.Get()
	data, err := client.Do("GET", key)
	if err != nil {
		return nil, err
	}

	client.Close()
	return data, nil
}

func Del(key string) error {
	client := pool.Get()
	_, err := client.Do("DEL", key)
	if err != nil {
		return err
	}

	client.Close()
	return nil
}
