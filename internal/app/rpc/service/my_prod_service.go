//
// Package services
// @Description：
// @Author：liuzezhong 2021/6/28 4:27 下午
// @Company cloud-ark.com
//
package services

import (
	"context"
)

type ProductService struct {
}

func (service *ProductService) ProService(_ context.Context, request *ProdRequest) (*ProdResponse, error) {
	return &ProdResponse{
		ProStock: request.ProId,
	}, nil
}
