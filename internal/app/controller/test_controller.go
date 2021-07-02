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
	"micro-project/internal/pkg/db/mysql"
)

var services = service_impl.NewTestService()

var balances balance.ILoadBalance = balance.RandomBalance{}

//var balances balance.ILoadBalance = balance.HashBalance{}
//var balances balance.ILoadBalance = &balance.RoundRobinBalance{}
//var balances balance.ILoadBalance = &balance.WeightRoundRobinBalance{}

type Test struct {
	Age int `json:"age"`
}

type T1 struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type TestController struct {
	Ctx irisContext.Context
}

func (c TestController) GetName() string {
	var conn = mysql.GetConn()
	fmt.Println(conn)
	var t1 = T1{}
	conn.Table("t1").First(&t1, "a=1")

	var t1Data = &T1{
		A: 1,
		B: "测试两下",
	}
	conn.Create(t1Data)
	var datas = []T1{
		T1{A: 2, B: "bbb"},
		T1{A: 3, B: "bbb1"},
		T1{A: 4, B: "bbb2"},
		T1{A: 5, B: "bbb3"},
		T1{A: 6, B: "bbb4"},
	}
	conn.Create(&datas)

	fmt.Println(t1)
	defer c.Ctx.Next()
	//var userServices = discover.DiscoverServices("userServic1")
	//var userService, err = balances.DoBalance(userServices, "userServic1", "10.0.10.253")
	//if err != nil {
	//	fmt.Println("出现错误，", err.Error())
	//} else {
	//	fmt.Println(userService.GetUrl())
	//}

	//for _, service := range userService {
	//	//fmt.Println(service)
	//	var httpUrl = "http://" + service.Host + ":" + strconv.Itoa(service.Port) + "/test/name"
	//
	//	client := &http.Client{}
	//	req, _ := http.NewRequest("GET", httpUrl, nil)
	//	req.Header = c.Ctx.Request().Header
	//	req.Header.Set("Uber-Trace-Id", c.Ctx.Request().Header.Get("Uber-Trace-Id"))
	//
	//	// 发送请求
	//	resp, _ := client.Do(req)
	//	defer resp.Body.Close()
	//	body, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))
	//}

	var name = services.GetName("陆小凤")
	fmt.Println(name)
	return name
}

func (c TestController) GetAge() Test {
	var age = services.GetAge("18")
	fmt.Println(age)
	defer c.Ctx.Next()
	return Test{
		Age: age,
	}
}

func (c TestController) GetAbc() []Test {
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
	return result
}