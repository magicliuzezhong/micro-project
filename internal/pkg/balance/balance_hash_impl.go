//
// Package balance
// @Description：hash负载均衡实现
// @Author：liuzezhong 2021/6/30 5:36 下午
// @Company cloud-ark.com
//
package balance

import (
	"errors"
	"fmt"
	"hash/crc32"
	"micro-project/internal/pkg/common"
	"micro-project/internal/pkg/util"
	"sync"
)

//
// HashBalance
// @Description: hash
//
type hashBalance struct{}

//
// @Description: 用于单例模式
//
var hashBalanceOnce sync.Once

//
// @Description: hashBalance实例
//
var hashBalanceInstance *hashBalance

//
// NewHashBalance
// @Description: 创建hashBalance实例
// @return ILoadBalance 负载均衡接口
//
func NewHashBalance() ILoadBalance {
	hashBalanceOnce.Do(func() {
		hashBalanceInstance = &hashBalance{}
	})
	return hashBalanceInstance
}

//
// DoBalance
// @Description: hash负载均衡实现（根据本地ip进行hash取模，并且映射到具体的实例中）
// @receiver r HashBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (r hashBalance) DoBalance(instances []*common.ServiceInstance, _ ...string) (*common.ServiceInstance, error) {
	var localIp = util.GetLocalIp()
	if localIp == "" {
		return nil, errors.New("获取本地ip失败")
	}

	var lens = len(instances)
	if lens == 0 {
		return nil, fmt.Errorf("实例找不到")
	}
	crcTable := crc32.MakeTable(crc32.IEEE)
	hashVal := crc32.Checksum([]byte(localIp), crcTable)
	index := int(hashVal) % lens
	var instance = instances[index]
	instance.CallTimes++
	return instance, nil
}
