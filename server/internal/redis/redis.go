package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"taskFive/server/config"
	e "taskFive/server/lib/err"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg config.Config) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &Redis{client: rdb}
}

func (r *Redis) SetManyValue(pairs []interface{}) error {
	var values []interface{}
	for i, v := range pairs {
		if i%2 == 1 {
			values = append(values, v)
		}
	}
	if len(values) == 0 {
		return nil
	}
	err := r.client.MSet(context.Background(), pairs...).Err()
	if err != nil {
		return err
	}
	for _, v := range values {
		key := fmt.Sprintf("%v", v)
		err := r.client.Expire(context.Background(), key, time.Minute).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Redis) SetValue(key string, value []byte, t time.Duration) error {
	err := r.client.Set(context.Background(), key, value, t).Err()
	if err != nil {
		log.Print("error", err)
		return err
	}

	return nil
}

func (r *Redis) DeleteValue(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}

	defer func() { err = e.WrapIfErr("can`t delete from redis: ", err) }()

	return nil
}

func (r *Redis) GetAllValue() ([]interface{}, []string, error) {
	keys, err := r.client.Keys(context.Background(), "*").Result()
	if err != nil {
		return nil, nil, err
	}

	values, err := r.client.MGet(context.Background(), keys...).Result()
	if err != nil {
		return values, nil, err
	}

	defer func() { err = e.WrapIfErr("can`t get values from redis: ", err) }()
	return values, keys, nil
}
