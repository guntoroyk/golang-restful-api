package service

import (
	"context"
	"github.com/guntoroyk/golang-restful-api/model/web"
)

type CategoryService interface {
	Create(ctx context.Context, request web.CategoryCreateRequest) *web.CategoryResponse
	Update(ctx context.Context, request web.CategoryUpdateRequest) (*web.CategoryResponse, error)
	Delete(ctx context.Context, categoryId int) error
	FindById(ctx context.Context, categoryId int) (*web.CategoryResponse, error)
	FindAll(ctx context.Context) []*web.CategoryResponse
}
