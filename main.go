package main

import (
	"os"
	"wadary/apis"
	"wadary/database/redis"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	redisCfg := redis.RedisConfig{
		Address:  os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       os.Getenv("REDIS_DB"),
		ENV:      os.Getenv("REDIS_ENV"),
	}
	cache, redisClient := redis.NewRedisCache(redisCfg)
	defer redisClient.Close()

	router := apis.NewRouter(port, cache)
	err := router.Start()
	if err != nil {
		panic(err)
	}
}
