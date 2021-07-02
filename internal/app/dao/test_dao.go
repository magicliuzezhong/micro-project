//
// Package dao
// @Description：
// @Author：liuzezhong 2021/6/25 6:48 下午
// @Company cloud-ark.com
//
package dao

type ITestDao interface {
	GetName(id string) string
	GetAge(id string) int
}
