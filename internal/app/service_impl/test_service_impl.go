//
// Package service_impl
// @Description：
// @Author：liuzezhong 2021/6/25 6:46 下午
// @Company cloud-ark.com
//
package service_impl

import (
	"micro-project/internal/app/dao_impl"
	"micro-project/internal/app/service"
	"sync"
)

var testDao = dao_impl.NewTestDao()

type testService struct {
}

var testServiceOnce sync.Once

var myService *testService

func NewTestService() service.ITestService {
	testServiceOnce.Do(func() {
		myService = &testService{}
	})
	return myService
}

func (service testService) GetName(id string) string {
	return testDao.GetName(id)
}

func (service testService) GetAge(id string) int {
	return testDao.GetAge(id)
}
