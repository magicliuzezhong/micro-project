//
// Package util
// @Description：配置文件工具类
// @Author：liuzezhong 2021/6/28 11:42 上午
// @Company cloud-ark.com
//
package util

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

//
// Application
// @Description: application配置文件
//
type Application struct {
	Server struct { //服务
		Port       int    `yaml:"port"`
		ServerName string `yaml:"server_name"`
	} `yaml:"server"`
	DiscoverConf struct {
		ConsulConf struct {
			Ip     string   `yaml:"ip"`
			Port   int      `yaml:"port"`
			Tag    []string `yaml:"tag"`
			Weight string   `yaml:"weight"`
		} `yaml:"consul"`
	} `yaml:"discover"`
	JeagerConf struct {
		Enabled bool   `yaml:"enabled"`
		Url     string `yaml:"url"`
		Type    string `yaml:"type"`
		Param   int    `yaml:"param"`
	} `yaml:"jeager"`
	DBInfo struct {
		Mysql struct {
			Url         string `yaml:"url"`
			Username    string `yaml:"username"`
			Password    string `yaml:"password"`
			Schema      string `yaml:"schema"`
			MaxIdleConn int    `yaml:"max_idle_conn"`
			MaxOpenConn int    `yaml:"max_open_conn"`
			LogPath     string `yaml:"log_path"`
		}
		Redis struct {
			RedisAlone struct {
				Enabled  bool   `yaml:"enabled"`
				Url      string `yaml:"url"`
				Port     int    `yaml:"port"`
				Password string `yaml:"password"`
			} `yaml:"redis_alone"`
			RedisCluster struct {
				Enabled          bool               `yaml:"enabled"`
				Password         string             `yaml:"password"`
				RedisClusterInfo []RedisClusterInfo `yaml:"redis_cluster_info"`
			} `yaml:"redis_cluster"`
		} `yaml:"redis"`
	} `yaml:"db"`
}

type RedisClusterInfo struct {
	Url  string `yaml:"url"`
	Port int    `yaml:"port"`
}

var application Application

var readConfigOnce sync.Once

//
// GetApplication
// @Description: 获取配置信息
// @return Application 应用
//
func GetApplication() Application {
	readConfigOnce.Do(func() {
		// 加载文件
		yamlFile, err := ioutil.ReadFile("./configs/application.yaml")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(yamlFile))
		// 将读取的yaml文件解析为响应的 struct
		err = yaml.Unmarshal(yamlFile, &application)
		if err != nil {
			panic("读取application.yaml文件异常" + err.Error())
		}
	})
	return application
}
