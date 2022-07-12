package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// redis

// 定义一个全局变量
var redisdb *redis.Client

const keepmitu = time.Minute * 10

func Init() {
	err := Open()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n", err)
		//panic(err)
	}
	defer Close()
	// 存普通string类型，10分钟过期
	redisdb.Set("test:name", "科科儿子", time.Minute*10)
	// 存hash数据
	redisdb.HSet("test:class", "521", 42)
	// 存list数据
	redisdb.RPush("test:list", 1) // 向右边添加元素
	redisdb.LPush("test:list", 2) // 向左边添加元素
	// 存set数据
	redisdb.SAdd("test:set", "apple")
	redisdb.SAdd("test:set", "pear")
	return
}
func Open() (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // redis地址
		Password: "",               // redis没密码，没有设置，则留空
		DB:       0,                // 使用默认数据库
	})
	_, err = redisdb.Ping().Result()
	return
}

func Close() {
	redisdb.Close()
}

// func Get(key string) (string, error) {
// 	err := Open()
// 	if err != nil {
// 		fmt.Printf("connect redis failed! err : %v\n", err)
// 		//panic(err)
// 	}
// 	defer Close()

// 	sType, err := redisdb.Type(key).Result()
// 	if err != nil {
// 		fmt.Println("get type failed :", err)
// 	}
// 	switch sType {
// 	case "string":
// 		break
// 	case "set":
// 		break
// 	case "hash":
// 		break
// 	case "list":
// 		break
// 	case "zset":
// 		break
// 	}

// 	value, err := redisdb.Get(key).Result()
// 	if err != nil {
// 		fmt.Printf("credisdb.Get(key).Result()! err : %v\n", err)
// 		//panic(err)
// 	}
// 	return value, err
// }

func Get(key string) (interface{}, error) {
	err := Open()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n", err)
		//panic(err)
	}
	defer Close()

	sType, err := redisdb.Type(key).Result()
	if err != nil {
		fmt.Println("get type failed :", err)
	}
	var val interface{}
	switch sType {
	case "string":
		value, er := redisdb.Get(key).Result()
		if err != nil {
			fmt.Println("credisdb.Get(key).Result()! err : ", err)
			//panic(err)
			err = er
		}
		val = value
		break
	case "set":
		break
	case "hash":
		break
	case "list":
		break
	case "zset":
		break
	default:
		break
	}

	return val, err
}

// func Set(key string, value string) {
// 	err := Open()
// 	if err != nil {
// 		fmt.Printf("connect redis failed! err : %v\n", err)
// 		panic(err)
// 	}
// 	defer Close()

// 	redisdb.Set(key, value, keepmitu)
// 	if err != nil {
// 		fmt.Printf("credisdb.Get(key).Result()! err : %v\n", err)
// 		panic(err)
// 	}
// 	return
// }

func Set(paras ...interface{}) {
	err := Open()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n", err)
		panic(err)
	}
	defer Close()

	key, _ := paras[0].(string)
	var dt time.Duration
	if len(paras) > 2 {
		dt, _ = paras[2].(time.Duration)
	} else {
		dt = keepmitu
	}
	redisdb.Set(key, paras[1], dt)
	if err != nil {
		fmt.Printf("credisdb.Get(key).Result()! err : %v\n", err)
		panic(err)
	}
	return
}

func Clear() {

	err := Open()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n", err)
		panic(err)
	}
	defer Close()

	var cursor uint64
	keys, cursor, err := redisdb.Scan(cursor, "*", 100).Result()
	if err != nil {
		fmt.Println("scan keys failed err:", err)
		panic(err)
	}
	for _, key := range keys {

		fmt.Println("key:", key)
		sType, err := redisdb.Type(key).Result()
		if err != nil {
			fmt.Println("get type failed :", err)
			return
		}
		fmt.Println("keyType:", sType)
		redisdb.Do("del", key)
		//fmt.Println("key:",key)
		// sType,err := redisdb.Type(ctx,key).Result()
		// if err !=nil{
		// 	fmt.Println("get type failed :",err)
		// 	return
		// }
		// fmt.Printf("key :%v ,type is %v\n",key,sType)
		// if sType == "string" {
		// 	val,err := redisdb.Get(ctx,key).Result()
		// 	if err != nil{
		// 		fmt.Println("get key values failed err:",err)
		// 		return
		// 	}
		// 	fmt.Printf("key :%v ,value :%v\n",key,val)
		// }else if sType == "list"{
		// 	val,err := redisdb.LPop(ctx,key).Result()
		// 	if err !=nil{
		// 		fmt.Println("get list value failed :",err)
		// 		return
		// 	}
		// 	fmt.Printf("key:%v value:%v\n",key,val)
		// }
	}
}
