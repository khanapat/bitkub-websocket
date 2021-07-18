package redis

import "github.com/gomodule/redigo/redis"

type SetDataNoExpireRedisFn func(key string, value interface{}) error

func NewSetDataNoExpireRedisFn(pool *redis.Pool) SetDataNoExpireRedisFn {
	return func(key string, value interface{}) error {
		conn := pool.Get()
		defer conn.Close()

		_, err := conn.Do("SET", key, value)
		if err != nil {
			return err
		}
		return nil
	}
}
