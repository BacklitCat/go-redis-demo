package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

// NewDefaultClient 根据redis配置初始化一个客户端
func NewDefaultClient() (redisClient *redis.Client, err error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // redis地址
		Password: "",               // redis密码，没有则留空
		DB:       0,                // 默认数据库，默认是0
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	if _, err = redisClient.Ping().Result(); err != nil {
		return nil, err
	}
	return redisClient, nil
}

func InitDemo() {
	redisClient, err := NewDefaultClient()
	if err != nil {
		//redis连接错误
		panic(err)
	}
	fmt.Println("Redis连接成功", redisClient)
}

func SetGetDemo() {
	redisClient, err := NewDefaultClient()
	if err != nil {
		panic(err)
	}
	//redisClient.Set("new_key", "new_value", time.Second*10)
	err = redisClient.Set("new_key", "new_value", 0).Err()
	if err != nil {
		panic(err)
	}
	value, err := redisClient.Get("new_key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

func SetNXDemo() {
	redisClient, err := NewDefaultClient()
	if err != nil {
		panic(err)
	}
	//redisClient.Set("new_key", "new_value", time.Second*10)
	err = redisClient.Set("new_key", "new_value", 0).Err()
	if err != nil {
		panic(err)
	}
	value, err := redisClient.Get("new_key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
