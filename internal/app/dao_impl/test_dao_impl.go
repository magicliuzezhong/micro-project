//
// Package dao_impl
// @Description：
// @Author：liuzezhong 2021/6/25 6:49 下午
// @Company cloud-ark.com
//
package dao_impl

import (
	"fmt"
	dao2 "micro-project/internal/app/dao"
	"micro-project/internal/pkg/db/mysql"
	"sync"
)

type testDao struct {
}

var testDaoOnce sync.Once

var dao *testDao

type T1 struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func NewTestDao() dao2.ITestDao {
	testDaoOnce.Do(func() {
		dao = &testDao{}
	})
	return dao
}

func (service testDao) GetName(id string) string {
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

	return "叶孤城" + id
}

func (service testDao) GetAge(id string) int {
	return 18
}
