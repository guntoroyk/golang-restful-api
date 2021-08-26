package repository

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/mock"
	"golang-restful-api/model/domain"
)

type MockCategoryRepositoryImpl struct {
	mock.Mock
}

func (mockRepository *MockCategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	args := mockRepository.Called()
	result := args.Get(0)
	return result.(domain.Category)
}

func (mockRepository *MockCategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	panic("implement me")
}

func (mockRepository *MockCategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category domain.Category) {
	panic("implement me")
}

func (mockRepository *MockCategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int) (domain.Category, error) {
	panic("implement me")
}

func (mockRepository *MockCategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	args := mockRepository.Called()
	result := args.Get(0)
	return result.([]domain.Category)
}
