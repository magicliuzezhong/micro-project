//
// Package balance
// @Description：随机负载均衡
// @Author：liuzezhong 2021/6/30 5:23 下午
// @Company cloud-ark.com
//
package balance

import (
	"errors"
	"math/rand"
	"micro-project/internal/pkg/common"
)

//
// RandomBalance
// @Description: 随机负载均衡实现
//
type RandomBalance struct {
}

//
// DoBalance
// @Description: 随机负载均衡
// @receiver r RandomLoadBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (r RandomBalance) DoBalance(instances []*common.ServiceInstance, _ ...string) (*common.ServiceInstance, error) {
	lens := len(instances)
	if lens == 0 {
		return nil, errors.New("实例找不到")
	}
	index := rand.Intn(lens)
	instance := instances[index]
	instance.CallTimes++
	return instance, nil
}
