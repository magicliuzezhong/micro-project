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

type Application struct {
	Server     Server     `yaml:"server"`
	ConsulConf ConsulConf `yaml:"consul"`
	JeagerConf JeagerConf `yaml:"jeager"`
	DBInfo     DBInfo     `yaml:"db"`
}

type Server struct {
	Port       string `yaml:"port"`
	ServerName string `yaml:"server_name"`
}

type ConsulConf struct {
	Ip     string
	Port   string
	Tag    string
	Weight string
}

//
// JeagerConf
// @Description: jeager配置
//
type JeagerConf struct {
	Enabled string `yaml:"enabled"`
	Url     string `yaml:"url"`
	Type    string `yaml:"type"`
	Param   string `yaml:"param"`
}

//
// DBInfo
// @Description: DB信息，其中包含了mysql、redis等等
//
type DBInfo struct {
	Mysql struct {
		Url         string `yaml:"url"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Schema      string `yaml:"schema"`
		MaxIdleConn string `yaml:"max_idle_conn"`
		MaxOpenConn string `yaml:"max_open_conn"`
		LogPath     string `yaml:"log_path"`
	}
	redis struct {
		Url      string `yaml:"url"`
		Password string `yaml:"password"`
	}
}

var application Application

var readConfigOnce sync.Once

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