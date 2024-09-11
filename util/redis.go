package util

import (
	"github.com/go-redis/redis"
	"github.com/randyhg/test-log-scanner/config"
	"github.com/randyhg/test-log-scanner/util/mylog"
	"sync"
)

var redisCache *FwRedis
var onceRedis sync.Once

type FwRedis struct {
	client redis.UniversalClient
}

func GetRedisCache() redis.UniversalClient {
	return redisCache.client
}

func InitRedis() {
	onceRedis.Do(func() {
		redisCache = new(FwRedis)
		redisCache.connectDB(config.Instance.RedisCache)
	})
}

func (r *FwRedis) connectDB(conf config.RedisConfig) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Host,
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		mylog.Fatal("redis connect ping failed, err:", err)
	} else {
		mylog.Debug("redis connect ping response:", "pong", pong)
		r.client = client
	}
}
