package main

import (
	"context"
	"fmt"
	_ "github.com/LiYunF/tool/logger"
	"github.com/LiYunF/tool/pool"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

func main() {
	//fmt.Print("ok")
	//logger.InitLogger("./log", "test.log")
	//a := 10
	////logger.L.Error("hh", errors.New("asd"), 0, a)
	//logger.L.Err(myErr.CreateError(20000, "hhh", a), 0)

	key := "user_linkage"
	if err := pool.InitRedisConnect("localhost:6379", "123456"); err != nil {
		panic(err)
	}
	if err := pool.InitRedisDataByMysql("root", "sql123", "localhost",
		"user_linkage", "ippool", "ip", key); err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	limit := rate.Every(time.Second) // 每秒一个
	limiter := rate.NewLimiter(limit, 10)
	q := time.Now()
	fmt.Println(q)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		_ = limiter.Wait(context.TODO())
		go func() {
			defer wg.Done()

			_, err := pool.GetTop1Ip(key)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Millisecond * (3000))

		}()
	}
	wg.Wait()
	q = time.Now()
	fmt.Println(q)

	//res, err := pool.GetTop1Ip(key)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(res)
}
