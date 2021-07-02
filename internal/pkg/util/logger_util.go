//
// Package util
// @Description：日志工具类
// @Author：liuzezhong 2021/6/25 5:43 下午
// @Company cloud-ark.com
//
package util

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

//最大日志文件大小 2M
var maxLogLength int64 = 1024 * 1024 * 10

//通用日志，用于记录程序通用日志信息
var commonLogger *logrus.Logger

//
// init
// @Description: 程序初始化
//
func init() {
	initCommonLogger()
}

//
// GetLogger
// @Description: 获取普通日志
// @return *logrus.Logger 日志
//
func GetLogger() *logrus.Logger {
	return commonLogger
}

//
// initCommonLogger
// @Description: 初始化日志收集系统
//
func initCommonLogger() {
	var content, err = rotatelogs.New("./log/common/sys.log"+"-%Y%m%d%H%M",
		rotatelogs.WithLinkName("./log/common/sys.log"), // 生成软链，指向最新日志文件
		//MaxAge and RotationCount cannot be both set  两者不能同时设置
		//rotatelogs.WithMaxAge(time.Second * 10), //clear 最小分钟为单位
		rotatelogs.WithRotationSize(maxLogLength), //最大文件大小
		//rotatelogs.WithRotationCount(5),           //number 默认7份 大于7份 或到了清理时间 开始清理
		rotatelogs.WithRotationTime(time.Minute), //rotate 最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
	)
	if err != nil {
		logrus.Printf("创建rotatelogs失败: %s", err)
		return
	}
	commonLogger = logrus.New()
	commonLogger.SetOutput(io.MultiWriter(content, os.Stdout))
}
