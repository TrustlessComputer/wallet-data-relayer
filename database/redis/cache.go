package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type IRedisCache interface {
	GetAll() ([]string, error)
	Exists(key string) (*bool, error)
	SetData(key string, value interface{}) error
	SetStringData(key string, value string) error
	SetStringDataWithExpTime(key string, value string, exipredIn int) error
	GetData(key string) (*string, error)
	Delete(key string) error
	SetDataWithExpireTime(key string, value interface{}, exipredIn int) error //exipredIn second
}

type RedisConfig struct {
	Address  string
	Password string
	DB       string
	ENV      string
}

type redisCache struct {
	cfg    RedisConfig
	client *redis.Client
	ctx    context.Context
}

const REDIS_CACHE_EXPIRED_TIME int = 86400 //a day

func NewRedisCache(cfg RedisConfig) (*redisCache, *redis.Client) {
	r := new(redisCache)
	ctx := context.Background()
	redisDB, err := strconv.Atoi(cfg.DB)
	if err != nil {
		//panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:       redisDB,      // use default DB
	})

	r.cfg = cfg
	r.client = rdb
	r.ctx = ctx
	return r, rdb
}

func (r *redisCache) SetStringData(key string, value string) error {
	timeD := time.Duration(int32(REDIS_CACHE_EXPIRED_TIME)) * time.Second
	err := r.client.Set(key, value, timeD).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetStringDataWithExpTime(key string, value string, exipredIn int) error {
	timeD := time.Duration(int32(exipredIn)) * time.Second
	err := r.client.Set(key, value, timeD).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetData(key string, value interface{}) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(key, valueByte, time.Second*time.Duration(REDIS_CACHE_EXPIRED_TIME)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) SetDataWithExpireTime(key string, value interface{}, exipredIn int) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	timeD := time.Duration(int32(exipredIn)) * time.Second
	err = r.client.Set(key, valueByte, timeD).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) GetData(key string) (*string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (r *redisCache) Delete(key string) error {
	err := r.client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Exists(key string) (*bool, error) {
	value, err := r.client.Exists(key).Result()
	if err != nil {
		return nil, err
	}
	res := value > 0
	return &res, nil
}

func (r *redisCache) GetAll() ([]string, error) {
	var keys []string
	var err error
	//ctx := context.Background()
	c := 0
	keys, _, err = r.client.Scan(uint64(c), "*", 100000).Result()
	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		keys = append(keys, key)
	}

	return keys, err
}
