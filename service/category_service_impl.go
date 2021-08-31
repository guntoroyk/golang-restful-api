package service

import (
	"context"
	"errors"
	"github.com/guntoroyk/golang-restful-api/helper"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"github.com/guntoroyk/golang-restful-api/model/web"
	"github.com/guntoroyk/golang-restful-api/repository"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) *web.CategoryResponse {

	category := domain.Category{
		Name: request.Name,
	}

	category = service.CategoryRepository.Save(ctx, category)

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) (*web.CategoryResponse, error) {

	category, err := service.CategoryRepository.FindById(ctx, request.Id)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return nil, errors.New("category not found")
	}

	category.Name = request.Name

	category = service.CategoryRepository.Update(ctx, category)

	return helper.ToCategoryResponse(category), nil
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) error {

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return errors.New("category not found")
	}

	service.CategoryRepository.Delete(ctx, category)

	return nil
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) (*web.CategoryResponse, error) {

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return nil, err
	}

	return helper.ToCategoryResponse(category), nil
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []*web.CategoryResponse {

	categories := service.CategoryRepository.FindAll(ctx)

	return helper.ToCategoryResponses(categories)
}
