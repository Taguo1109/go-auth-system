package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

/**
 * @File: redis.go.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/10 上午10:19
 * @Software: GoLand
 * @Version:  1.0
 */

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("❌ Redis 連線失敗：", err)
	}
	log.Println("✅ Redis 已連線")
}
