package redis

import (
	"axis/ecommerce-backend/configs"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Client struct {
	client *redis.Client
}

type RememberFunc func(f interface{}) (string, error)

var ctx = context.Background()

func NewRedisClient() *Client {
	opts := configs.RedisOptions()
	if opts == nil {
		log.Println("nil setup")
		return nil
	}
	rdb := redis.NewClient(opts)
	return &Client{
		client: rdb,
	}
}

func (c *Client) GetRedis() *redis.Client {
	return c.client
}

func (c *Client) Set(key string, value interface{}) error {
	return c.store(key, value, 0)
}

func (c *Client) Get(key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("key does not exist")
		}
		return "", err
	}
	return val, nil
}

func (c *Client) Remember(key string, def RememberFunc, param interface{}) (string, error) {
	v, err := c.Get(key)
	if err != nil {
		u, err := def(param)
		if err != nil {
			return "", err
		}

		err = c.Set(key, u)
		if err != nil {
			return "", err
		}
		return c.Get(key)
	}
	return v, err
}

func (c *Client) store(k string, v interface{}, t time.Duration) error {
	err := c.client.Set(ctx, k, v, t).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Del(keys ...string) (int64, error) {
	tt, err := c.client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, err
	}
	return tt, nil
}

func (c *Client) Ping() error {
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		//log.Println(err)
		return err
	}
	log.Println("Redis Pong")
	return nil
}

func (c *Client) SetWithTime(key string, value interface{}, duration time.Duration) error {
	return c.store(key, value, duration)
}

func (c *Client) RememberWithTime(key string, def RememberFunc, param interface{}, d time.Duration) (string, error) {
	v, err := c.Get(key)
	if err != nil {
		u, err := def(param)
		if err != nil {
			return "", err
		}

		err = c.SetWithTime(key, u, d)
		if err != nil {
			return "", err
		}
		return c.Get(key)
	}
	return v, err
}
