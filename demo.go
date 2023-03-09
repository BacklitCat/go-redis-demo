package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"reflect"
	"time"
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

// SetNXDemo
//第1次执行SetNX状态： true
//第2次执行SetNX状态： false
//new_nx_value

func SetNXDemo() {
	redisClient, err := NewDefaultClient()
	if err != nil {
		panic(err)
	}
	status, err := redisClient.SetNX("new_nx_key", "new_nx_value", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("第1次执行SetNX状态：", status)
	status, err = redisClient.SetNX("new_nx_key", "new_nx_value", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("第2次执行SetNX状态：", status)
	value, err := redisClient.Get("new_nx_key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}

// GetSetDemo
//第1次执行GetSet，没有旧值会抛出错误： redis: nil
//执行Set，设置旧值
//第2次执行GetSet，返回旧值： new_get_set_value1
//执行Get，得到新值： new_get_set_value2
//新值的TTL： -1s
//重新设置TTL避免Demo弄脏Redis，执行： expire new_get_set_key 10: true

func GetSetDemo() {

	redisClient, err := NewDefaultClient()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.GetSet("new_get_set_key", "new_get_set_value1").Result()
	if err != nil {
		fmt.Println("第1次执行GetSet，没有旧值会抛出错误：", err)
		_, err = redisClient.Set("new_get_set_key", "new_get_set_value1", 10*time.Second).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("执行Set，设置旧值")
	}

	//在使用 DEL、SET、GETSET 等命令会清除对应的key的过期时间。
	val, err = redisClient.GetSet("new_get_set_key", "new_get_set_value2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("第2次执行GetSet，返回旧值：", val)

	val, err = redisClient.Get("new_get_set_key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("执行Get，得到新值：", val)

	ttlTime, err := redisClient.TTL("new_get_set_key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("新值的TTL：", ttlTime)

	status := redisClient.Expire("new_get_set_key", 10*time.Second)
	fmt.Println("重新设置TTL避免Demo弄脏Redis，执行：", status)
}
func MGetMSet() {
	redisClient, err := NewDefaultClient()
	if err != nil {
		panic(err)
	}

	_, err = redisClient.MSet("k1", "v1", "k2", "v2", "k3", "v3").Result()
	if err != nil {
		panic(err)
	}

	values, err := redisClient.MGet("k1", "k2", "k3").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(values, reflect.TypeOf(values)) //[v1 v2 v3] []interface {}
	delNum, err := redisClient.Del("k1", "k2", "k3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(delNum) // 3

}

// BLPOPDemo
//4
//v1
//[list v2]
//[list v3]
//[list v4]
//(wait 5s)
//Process finished with the exit code 0

func BLPOPDemo() {
	redisClient, err := NewDefaultClient()
	if err != nil {
		panic(err)
	}
	result, err := redisClient.RPush("list", "v1", "v2", "v3", "v4").Result()
	if err != nil {
		return
	}
	fmt.Println(result)

	str, err := redisClient.LPop("list").Result()
	if err != nil {
		return
	}
	fmt.Println(str)

	for i := 0; i < 6; i++ {
		strs, err := redisClient.BLPop(5*time.Second, "list").Result()
		if err != nil {
			return
		}
		fmt.Println(strs)
	}
}
