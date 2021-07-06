//
// Package balance
// @Description：加权轮询
// @Author：liuzezhong 2021/6/30 5:50 下午
// @Company cloud-ark.com
//
package balance

import (
	"fmt"
	"micro-project/internal/pkg/common"
	"sync"
)

//
// weightRoundRobinBalance
// @Description: 加权轮询
//
type weightRoundRobinBalance struct {
	Index  int64
	Weight int64
}

//
// @Description: 单例
//
var weightRoundRobinBalanceOnce sync.Once

//
// @Description: 实例
//
var weightRoundRobinBalanceInstance *weightRoundRobinBalance

//
// NewRoundRobinWeightBalance
// @Description: 获取加权轮询实例
// @return ILoadBalance 负载均衡接口
//
func NewRoundRobinWeightBalance() ILoadBalance {
	weightRoundRobinBalanceOnce.Do(func() {
		weightRoundRobinBalanceInstance = &weightRoundRobinBalance{}
	})
	return weightRoundRobinBalanceInstance
}

//
// DoBalance
// @Description: 加权轮询
// @receiver r weightRoundRobinBalance
// @param instances 实例
// @return *common.ServiceInstance 返回实例
// @return error 错误
//
func (r *weightRoundRobinBalance) DoBalance(instances []*common.ServiceInstance,
	_ ...string) (*common.ServiceInstance, error) {
	lens := len(instances)
	if lens == 0 {
		return nil, fmt.Errorf("实例找不到")
	}
	var instance = r.getInstance(instances)
	instance.CallTimes++
	return instance, nil
}

//
// getInstance
// @Description: 获取实例
// @receiver r weightRoundRobinBalance
// @param instances 实例
// @return *common.ServiceInstance 实例
//
func (r *weightRoundRobinBalance) getInstance(instances []*common.ServiceInstance) *common.ServiceInstance {
	gcd := getGCD(instances)
	for {
		r.Index = (r.Index + 1) % int64(len(instances))
		if r.Index == 0 {
			r.Weight = r.Weight - gcd
			if r.Weight <= 0 {
				r.Weight = getMaxWeight(instances)
				if r.Weight == 0 {
					return &common.ServiceInstance{}
				}
			}
		}
		if instances[r.Index].Weight >= int(r.Weight) {
			return instances[r.Index]
		}
	}
}

//
// gcd
// @Description: 计算两个数的最大公约数
// @param a a
// @param b b
// @return int64 最大公约数
//
func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

//
// getGCD
// @Description: 计算多个数的最大公约数
// @param instances 实例
// @return int64 最大公约数
//
func getGCD(instances []*common.ServiceInstance) int64 {
	var weights []int64

	for _, instance := range instances {
		weights = append(weights, int64(instance.Weight))
	}

	var g = weights[0]
	for i := 1; i < len(weights)-1; i++ {
		oldGcd := g
		g = gcd(oldGcd, weights[i])
	}
	return g
}

//
// getMaxWeight
// @Description: 获取最大权重
// @param instances 实例
// @return int64 最大公约数
//
func getMaxWeight(instances []*common.ServiceInstance) int64 {
	var max int64 = 0
	for _, instance := range instances {
		if int64(instance.Weight) >= max {
			max = int64(instance.Weight)
		}
	}
	return max
}
