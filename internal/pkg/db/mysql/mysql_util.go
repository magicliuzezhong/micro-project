//
// Package mysql
// @Description：mysql工具类（注：由于规范问题，请勿在非dao包里获取mysql连接）
// @Author：liuzezhong 2021/7/1 6:16 下午
// @Company cloud-ark.com
//
package mysql

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"micro-project/internal/pkg/util"
	"os"
	"sync"
	"time"
)

var dbs *gorm.DB

// mysql初始化错误信息
var dbInitError error

//暂停管道
var waitSendChannel = make(chan struct{}, 1)

//恢复管道
var resumeSendChannel = make(chan struct{}, 1)

// 初始化mysql独占锁，主要用于防止多链路切换时切换异常的发生
var initMysqlMutex sync.Mutex

//
// init
// @Description: 程序初始化
//
func init() {
	go mysqlExplore() //开启探活协程
	initMysqlDB()
}

//
// GetConn
// @Description: 获取mysql连接
// @return *sql.DB db
// @return error 错误
//
func GetConn() *gorm.DB {
	return dbs
}

//
// ResetMysqlDB
// @Description: 重置mysql数据库
//
func ResetMysqlDB() {
	initMysqlDB()
}

//
// mysqlExplore
// @Description: mysql探活程序
//
func mysqlExplore() {
	for {
		select {
		case <-waitSendChannel:
			<-resumeSendChannel
		case <-time.After(time.Second * 60):
			//此处需要进行加锁操作，因为如果mysql正在进行切换的时候dbs被置为nil了的话，此处会出现报错的问题
			initMysqlMutex.Lock()
			if dbs != nil {
				var dbSql, err = dbs.DB()
				if err != nil {
					initMysqlMutex.Unlock()
					continue
				}
				err = dbSql.Ping()
				if err == nil {
					util.GetLogger().Info("====mysql探活====")
				}
			}
			initMysqlMutex.Unlock()
		}
	}
}

//
// initMysqlDB
// @Description: 初始化mysqlDB
//
func initMysqlDB() {
	initMysqlMutex.Lock()
	defer initMysqlMutex.Unlock()
	if len(resumeSendChannel) == 1 {
		<-resumeSendChannel //如果没有等待了要先释放，防止出现刚暂停又恢复
	}
	if len(waitSendChannel) == 0 {
		waitSendChannel <- struct{}{} //程序阻塞
	}
	dbs = nil
	dbInitError = nil
	var mysqlInfo = util.GetApplication().DBInfo.Mysql
	if mysqlInfo.Url == "" { //url不存在，说明需要进行初始化
		if len(resumeSendChannel) == 1 {
			<-resumeSendChannel //如果没有等待了要先释放，防止出现刚暂停又恢复
		}
		if len(waitSendChannel) == 0 {
			waitSendChannel <- struct{}{} //程序阻塞
		}
		return
	}
	var defaultLogPath = "./log/common/mysql.log" //缺省的mysql日志地址
	if mysqlInfo.LogPath != "" {
		defaultLogPath = mysqlInfo.LogPath
	}
	var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlInfo.Username, mysqlInfo.Password, mysqlInfo.Url, mysqlInfo.Schema)
	content, err1 := rotatelogs.New(defaultLogPath+"-%Y%m%d%H%M",
		rotatelogs.WithLinkName(defaultLogPath),   // 生成软链，指向最新日志文件
		rotatelogs.WithRotationSize(1024*1024*10), //最大文件大小10M
	)
	if err1 != nil {
		log.Printf("创建rotatelogs失败: %s\n", err1)
	}
	newLogger := logger.New(
		log.New(io.MultiWriter(content, os.Stdout), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 禁用彩色打印，最好别使用这个，否则日志打印出来不太好看
		},
	)
	var db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		dbInitError = err
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		dbInitError = err
		return
	}
	var defaultMaxIdleConn = 10
	if mysqlInfo.MaxIdleConn > 0 {
		defaultMaxIdleConn = mysqlInfo.MaxIdleConn
	}
	var defaultMaxOpenConn = 100
	if mysqlInfo.MaxOpenConn > 0 {
		defaultMaxOpenConn = mysqlInfo.MaxOpenConn
	}
	sqlDB.SetMaxIdleConns(defaultMaxIdleConn) //设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(defaultMaxOpenConn) //设置打开数据库连接的最大数量。
	sqlDB.SetConnMaxLifetime(time.Hour)       //设置了连接可复用的最大时间。
	dbs = db
	dbInitError = nil //初始化错误信息设置为空，因为已经设置完成了并且没有出错
	if len(resumeSendChannel) == 0 {
		resumeSendChannel <- struct{}{} //恢复探活程序
		if len(waitSendChannel) == 1 {  //进行恢复后要将暂停逻辑取消掉，防止出现刚恢复又暂停的情况
			<-waitSendChannel
		}
	}
}
