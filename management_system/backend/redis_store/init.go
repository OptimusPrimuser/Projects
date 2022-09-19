package redis_store

import (
	"encoding/json"

	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client   *redis.Client
	ReClient *rejonson.Client
}

func (r *Redis) Init() {
	r.Client = redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
		},
	)
	r.ReClient = rejonson.ExtendClient(r.Client)
}

func (r *Redis) Close() {
	r.ReClient.Close()
}

func (r *Redis) Add(key string, data interface{}) error {
	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.ReClient.JsonSet(key, ".", string(byteData))
	return nil
}

func (r *Redis) Remove(key string) {
	r.ReClient.Del(key)
}

func (r *Redis) Get(key string) (string, error) {
	jsonString, err := r.ReClient.JsonGet(key).Result()
	if err != nil {
		return "", err
	}
	return jsonString, nil

}
