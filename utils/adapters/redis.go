package adapters

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr string
	DB   int
	URL  string
}

func CreateRedisClient(ctx context.Context, redisConfig RedisConfig) *redis.Client {
	var rdb *redis.Client
	if redisConfig.URL != "" {
		opt, _ := redis.ParseURL(redisConfig.URL)
		rdb = redis.NewClient(opt)
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: redisConfig.Addr,
			DB:   redisConfig.DB, // use default DB
		})
	}

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	return rdb
}

// Serialize converte um objeto em uma string JSON.
func Serialize(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Deserialize converte uma string JSON para o objeto indicado.
func Deserialize(dataStr string, target interface{}) error {
	return json.Unmarshal([]byte(dataStr), target)
}
