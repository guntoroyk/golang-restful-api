package cache

import (
	"context"
	"github.com/guntoroyk/golang-restful-api/model/domain"
)

type CategoryCache interface {
	SetCategory(ctx context.Context, category domain.Category) error
	SetCategoryBatch(ctx context.Context, categories []domain.Category) error
	Delete(ctx context.Context, category domain.Category) error
	GetCategory(ctx context.Context, categoryId int) (domain.Category, error)
	GetCategoryBatch(ctx context.Context) ([]domain.Category, error)
}
