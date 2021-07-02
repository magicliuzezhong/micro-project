//
// Package service
// @Description：
// @Author：liuzezhong 2021/6/25 6:46 下午
// @Company cloud-ark.com
//
package service

type ITestService interface {
	GetName(id string) string
	GetAge(id string) int
}
