// Package logger
// @Author lyf
// @Update lyf 2023.01
package logger

import (
	"context"
	"errors"
	"fmt"
	myErr "github.com/LiYunF/tool/myErr"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"runtime"
	"time"
)

//******************************************************************//
//							log 结构体								//
//******************************************************************//

// Log 封装以便日后更改，升级
// use slog
type Log struct {
	entity *slog.Logger
}

const (
	LevelDebug  = slog.LevelDebug
	LevelInfo   = slog.LevelInfo
	LevelWarn   = slog.LevelWarn
	LevelError  = slog.LevelError
	LevelDanger = slog.Level(12) //危险
	LevelPanic  = slog.Level(16) // 输出日志后panic
	LevelFatal  = slog.Level(20) // Fatal 致命错误，出现错误时程序无法正常运转，输出日志后程序退出
)

// Error
//
//	@description	logs at LevelError. If err is non-nil, Error appends Any(ErrorKey, err) to the list of attributes.
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Error(msg string, err error, deep int, param any) {
	l.entity.Error(msg, "err", err, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
}

// Err
//
//	@description	自定义error的Error级日志记录
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Err(myError *myErr.MyError, deep int) {
	l.entity.Error(fmt.Sprintf("code:%v", myError.Code), "err", myError,
		"param", fmt.Sprintf("{%#v}", myError.Data), "(source", l.getScour(deep+2)+")")
}

// Debug
//
//	@description	debug级别
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Debug(msg string, deep int, param any) {
	l.entity.Debug(msg, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
}

// Info
//
//	@description	info级别错误
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Info(msg string, deep int, param any) {
	l.entity.Info(msg, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
}

// Warn
//
//	@description	Warn级别错误
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Warn(msg string, deep int, param any) {
	l.entity.Warn(msg, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
}

// Fatal
//
//	@description	非必要不使用：致命错误，出现错误时程序无法正常运转，输出日志后程序退出(os.Exit(1))
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Fatal(msg string, deep int, param any) {
	l.entity.Log(context.Background(), LevelFatal, msg, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
	os.Exit(1)
}

// Danger
//
//	@description	非必要不使用：危险错误错误，出现错误时程序部分功能无法正常工作
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Danger(msg string, deep int, param any) {
	l.entity.Log(context.Background(), LevelDanger, msg, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
}

// Panic
//
//	@description	panic级别错误，输出日志后调用panic(msg)
//	@param	deep	函数栈深度，若调用log的位置是错误发生的位置，则输入0; 否则输入封装的深度
//	@param	param	错误时的相关参数信息，建议用fmt.Sprintf()
func (l *Log) Panic(msg string, deep int, param any) {
	l.entity.Log(context.Background(), LevelPanic, msg, "param", fmt.Sprintf("{%#v}", param), "(source", l.getScour(deep+2)+")")
	panic(msg)
}

// Partition
//
//	@description	打印分隔符============
func (l *Log) Partition(msg string) {
	l.entity.Info(msg, "Partition", "============================================================================================================")
}

//******************************************************************//
//							log 初始化								//
//******************************************************************//

var L *Log

func InitLogger(DirPath, NameFormat string) {
	if e := CreateLogger(DirPath, NameFormat); e != nil {
		panic("init logger error:" + e.Error())
	}
}

// CreateLogger 创建日志
// DirPath 文件夹名字 例如 ./log
// NameFormat 日志命名格式 例如 test.log
func CreateLogger(DirPath, NameFormat string) error {

	//init file
	dirPath := DirPath
	if !fileExists(dirPath) {
		if err := os.Mkdir(dirPath, 0777); err != nil {
			return err
		}
	}
	path := dirPath + "/" + NameFormat

	witter, err := getWriter(path)
	if err != nil {
		return err
	}

	//Text 类型
	handle := slog.HandlerOptions{

		Level: LevelDebug, //需要输出的日志级别，默认为info，测试环境可以写成debug

		//AddSource: true, //为true时，输出调用log的位置信息，但是一旦做了封装就只会输出封装的位置，所以建议为false

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				// 处理自定义级别
				level := a.Value.Any().(slog.Level)

				//自定义错误等级划分
				switch {

				case level < LevelInfo:
					a.Value = slog.StringValue("[ DEBUG ]")
				case level < LevelWarn:
					a.Value = slog.StringValue("[ INFO ]")
				case level < LevelError:
					a.Value = slog.StringValue("[ WARN ]")
				case level < LevelDanger:
					a.Value = slog.StringValue("[ ERROR ]")
				case level < LevelPanic:
					a.Value = slog.StringValue("[ DANGER ]")
				case level < LevelFatal:
					a.Value = slog.StringValue("[ PANIC ]")
				default:
					a.Value = slog.StringValue("[ FATAL ]")
				}
			}
			return a
		},
	}.NewTextHandler(witter)

	//赋值
	L = new(Log)
	L.entity = slog.New(handle)

	return nil
}

//******************************************************************//
//							util								//
//******************************************************************//

// getWriter 日志分割
func getWriter(path string) (io.Writer, error) {
	// 保存60天内的日志，每24小时(整点)分割一次日志
	return rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Hour*24*60),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

}

// fileExists 查看文件/文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// getScour 获取出问题的位置
func (l *Log) getScour(skip int) string {
	if skip < 2 {
		l.entity.Error("获取出错代码位置skip错误", errors.New("获取出错代码位置skip错误"), "skip", skip)
		return "源代码位置获取失败，skip错误"
	}
	pc, codePath, codeLine, ok := runtime.Caller(skip)
	if !ok {
		// 不ok，函数栈用尽了
		l.entity.Error("获取出错代码位置函数栈用尽", errors.New("获取出错代码位置函数栈用尽"), "skip", skip)
		return "源代码位置获取失败，函数栈用尽"
	}

	// 拼接文件名与所在行
	code := fmt.Sprintf("%s:%d func name:%s", codePath, codeLine, runtime.FuncForPC(pc).Name())
	return fmt.Sprintf(code)

}
