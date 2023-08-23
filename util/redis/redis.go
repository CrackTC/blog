package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"sora.zip/blog/config"
)

var rdb *redis.Client
var ctx context.Context

const Nil = redis.Nil

func Len(key string) (int64, error) {
	return rdb.LLen(ctx, key).Result()
}

func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func GetList(key string, start int64, count int64) ([]string, error) {
	cmd := rdb.Exists(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	if cmd.Val() == 0 {
		return nil, redis.Nil
	}
	return rdb.LRange(ctx, key, start, start+count-1).Result()
}

func Set(key string, value interface{}) error {
	expiration, err := time.ParseDuration(config.Get().RedisExpiration)
	if err != nil {
		log.Println("[WARN] Invalid redis expiration:", config.Get().RedisExpiration)
		expiration = time.Duration(0)
	}
	return rdb.Set(ctx, key, value, expiration).Err()
}

func SetList(key string, values []interface{}) error {
	if err := rdb.RPush(ctx, key, values...).Err(); err != nil {
		return err
	}
	expiration, err := time.ParseDuration(config.Get().BlogFetchInterval)
	if err != nil {
		log.Println("[WARN] Invalid redis expiration:", config.Get().BlogFetchInterval)
		return nil
	}
	return rdb.Expire(ctx, key, expiration).Err()
}

func init() {
	c := config.Get()
	opt, err := redis.ParseURL(c.RedisURL)
	if err != nil {
		log.Fatalln("[ERROR] Failed to parse redis url:", c.RedisURL)
	}

	rdb = redis.NewClient(opt)
	ctx = context.Background()
}
