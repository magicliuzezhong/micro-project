//
// Package common
// @Description：iris response返回model
// @Author：liuzezhong 2021/6/30 1:47 下午
// @Company cloud-ark.com
//
package common

//
// ResponseResult
// @Description: response结果
//
type ResponseResult struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
}
