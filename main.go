package main

import (
	"fmt"
	_ "github.com/LiYunF/tool/logger"
	"github.com/LiYunF/tool/pool"
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
	//if err := pool.InitRedisDataByMysql("root", "sql123", "localhost",
	//	"user_linkage", "ippool", "ip", "user_linkage"); err != nil {
	//	panic(err)
	//}
	res, err := pool.GetTop1Ip(key)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
