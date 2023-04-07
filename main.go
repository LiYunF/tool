package main

import (
	"fmt"
	"github.com/LiYunF/tool/logger"
	_ "github.com/LiYunF/tool/logger"
	myErr "github.com/LiYunF/tool/myErr"
)

func main() {
	fmt.Print("ok")
	logger.InitLogger("./log", "test.log")
	a := 10
	//logger.L.Error("hh", errors.New("asd"), 0, a)
	logger.L.Err(myErr.CreateError(20000, "hhh", a), 0)
}
