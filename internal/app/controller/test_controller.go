//
// Package controller
// @Description：测试控制器
// @Author：liuzezhong 2021/6/25 6:45 下午
// @Company cloud-ark.com
//
package controller

import (
	"fmt"
	irisContext "github.com/kataras/iris/v12/context"
	"micro-project/internal/app/service_impl"
	"micro-project/internal/pkg/balance"
	"micro-project/internal/pkg/common"
	"micro-project/internal/pkg/discover"
)

var services = service_impl.NewTestService()

//var balances = balance.NewConsistencyHashBalance()
//var balances = balance.NewHashBalance()
//var balances = balance.NewRandomBalance()
//var balances = balance.NewRoundRobinBalance()
var balances = balance.NewRoundRobinWeightBalance()

type Test struct {
	Age int `json:"age"`
}

type TestController struct {
	Ctx irisContext.Context
}

func (c TestController) GetName() {
	defer c.Ctx.Next()
	var userServices = discover.DiscoverServices("userService1")
	var userService, err = balances.DoBalance(userServices, "userService1", "10.0.10.253")
	if err != nil {
		fmt.Println("出现错误，", err.Error())
	} else {
		fmt.Println(userService.GetUrl())
		var httpUrl = userService.GetUrl() + "/test/name"
		fmt.Println("获取到的url：", httpUrl)
		//client := &http.Client{}
		//req, _ := http.NewRequest("GET", httpUrl, nil)
		//req.Header = c.Ctx.Request().Header
		//req.Header.Set("Uber-Trace-Id", c.Ctx.Request().Header.Get("Uber-Trace-Id"))
		//
		//// 发送请求
		//resp, _ := client.Do(req)
		//defer resp.Body.Close()
		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
	}
	panic("测试异常")

	//var name = services.GetName("陆小凤")
	//fmt.Println(name)
	//return name
	c.Ctx.Values().Set("val", "测试一下")
}

func (c TestController) GetAge() {
	var age = services.GetAge("18")
	defer c.Ctx.Next()
	c.Ctx.Values().Set("val", Test{
		Age: age,
	})

	c.Ctx.Values().Set("val", common.ResponseResult{
		Status: 200,
		Msg:    "我个人测试的请求",
		Data:   age,
	})
}

func (c TestController) GetBbb() string {
	return "西门吹雪"
}

func (c TestController) GetAbc() {
	var result = make([]Test, 0)
	result = append(result, Test{
		Age: 18,
	})
	result = append(result, Test{
		Age: 19,
	})
	result = append(result, Test{
		Age: 20,
	})
	defer c.Ctx.Next()
	c.Ctx.Values().Set("val", result)
}

func (c TestController) GetCba() {

}
