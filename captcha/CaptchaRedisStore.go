package captcha

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

// RedisStore 实现store接口类型
type RedisStore struct {
	Client     *redis.Client
	Expiration time.Duration
}

// Set 实现set方法，将数据写入Redis中，并设置好有效期
func (rs RedisStore) Set(id string, value string) error {
	return rs.Client.Set(id, value, rs.Expiration).Err()
}

// Verify 实现验证方法。从Redis中取出数据并进行比较，根据clear选择是否在验证后删除数据
func (rs RedisStore) Verify(id, answer string, clear bool) bool {
	v := rs.Get(id, clear)
	return v == answer
}

// Get 实现Redis数据的获取和删除
func (rs RedisStore) Get(id string, clear bool) string {
	val, err := rs.Client.Get(id).Result()
	if err != nil {
		return ""
	}
	if clear {
		if err = rs.Client.Del(id).Err(); err != nil {
			log.Println(err)
			return ""
		}
	}
	return val
}
