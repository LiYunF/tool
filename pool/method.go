package pool

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

var RedisClient *redis.Client
var QuarterIpNumber int
var Random *rand.Rand

// InitRedisConnect 初始化链接
func InitRedisConnect(addr, password string) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return errors.New("init redis connect error " + err.Error())
	}
	RedisClient = redisClient
	QuarterIpNumber = 50 //初始化前50个
	Random = rand.New(rand.NewSource(time.Now().UnixNano()))
	return nil
}

// InitRedisDataByMysql 初始化redis data
func InitRedisDataByMysql(username, password, host, database, tableName, column, redisKeyName string) error {
	if RedisClient == nil {
		return errors.New("you need to Init redis connect first")
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true&loc=Local",
		username, password, host, database)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return errors.New("connect mysql fail" + err.Error())
	}
	var ippool []string
	sql := fmt.Sprintf("select %v from %v", column, tableName)
	if err := db.Select(&ippool, sql); err != nil {
		return errors.New("select mysql fail" + err.Error())
	}
	//fmt.Println(ippool)
	QuarterIpNumber = 0
	for _, e := range ippool {
		if err := AddDataToZSet(redisKeyName, e, 0.0); err != nil {
			return errors.New("set data to zset fail:" + e + err.Error())
		}
		QuarterIpNumber += 1
	}
	//获取前1/4的数据
	QuarterIpNumber /= 4
	return nil
}
