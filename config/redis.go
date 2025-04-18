package config

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

/**
 * @File: redis.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/10 上午10:19
 * @Software: GoLand
 * @Version:  1.0
 */

var RDB *redis.Client
var Ctx = context.Background()

// InitRedis 初始化 Redis 客戶端，並加入重試機制以防 Redis 尚未準備好。
// 本機環境幾乎不會發生，但在 Docker compose 裡常見 Redis 還沒 ready 就連線的情況。
// Retry up to 10 times, waiting 2 seconds each time.
//
// 💡 記得 REDIS_ADDR 在 Docker 裡應設為 redis:6379（依照 compose 服務名稱）
//
// 可能錯誤範例：
// ❌ Redis connection error: dial tcp 172.xx.xx.xx:6379: connect: connection refused
func InitRedis() {

	rdb := redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_ADDR"),
		Username:  "default", // 固定
		Password:  os.Getenv("REDIS_PASS"),
		TLSConfig: &tls.Config{}, // TLS
		DB:        0,
	})

	for i := 1; i <= 10; i++ {
		_, err := rdb.Ping(Ctx).Result()
		if err == nil {
			RDB = rdb
			log.Println("✅ Redis 已連線")
			return
		}
		log.Printf("❌ Redis 連線失敗 (%d/10)：%v\n", i, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ Redis 最終連線失敗")
}
