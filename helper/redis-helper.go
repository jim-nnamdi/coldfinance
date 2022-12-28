package helper

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func SetRedisClient(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	beep, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("could not connect to redis: %s\n", beep)
		return
	}
	fmt.Println(beep)
	redisClient = client
}

func SetRedisData(ctx context.Context, key string, value interface{}) {
	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("cannot set key-value pair: %s %+v \n", key, value)
		return
	}
	redisClient.Incr(ctx, key)
}

func GetRedisData(ctx context.Context, key string) interface{} {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Printf("error fetching value: %+v", val)
		return nil
	}
	return val
}

func GetAllKeys(ctx context.Context, key string) []string {
	keys := []string{}
	iter := redisClient.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return keys
}

func GetSumOfValues(ctx context.Context, key string) uint {
	iter := redisClient.Scan(ctx, 0, key, 0).Iterator()
	sumOfValues := 0
	for iter.Next(ctx) {
		for v := range iter.Val() {
			log.Printf("%v\n", v)
			sumOfValues += v
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("err calc views : %s", err.Error())
		panic(err)
	}
	return uint(sumOfValues)
}

func GetIpAddress() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	IPs := make([]string, 0)
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IPs = append(IPs, ipnet.IP.To4().String())
			}
		}
	}
	fmt.Println(IPs[0])
	return IPs, nil
}
