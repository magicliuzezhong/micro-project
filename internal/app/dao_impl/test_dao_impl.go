//
// Package dao_impl
// @Description：
// @Author：liuzezhong 2021/6/25 6:49 下午
// @Company cloud-ark.com
//
package dao_impl

import (
	dao2 "micro-project/internal/app/dao"
	"sync"
)

type testDao struct {
}

var testDaoOnce sync.Once

var dao *testDao

func NewTestDao() dao2.ITestDao {
	testDaoOnce.Do(func() {
		dao = &testDao{}
	})
	return dao
}

func (service testDao) GetName(id string) string {
	return "叶孤城" + id
}

func (service testDao) GetAge(id string) int {
	return 18
}
