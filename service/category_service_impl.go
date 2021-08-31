package service

import (
	"context"
	"errors"
	"github.com/guntoroyk/golang-restful-api/cache"
	"github.com/guntoroyk/golang-restful-api/helper"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"github.com/guntoroyk/golang-restful-api/model/web"
	"github.com/guntoroyk/golang-restful-api/repository"
	"log"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	CategoryCache      cache.CategoryCache
}

func NewCategoryService(categoryRepository repository.CategoryRepository, categoryCache cache.CategoryCache) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		CategoryCache:      categoryCache,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) *web.CategoryResponse {

	category := domain.Category{
		Name: request.Name,
	}

	category = service.CategoryRepository.Save(ctx, category)

	categories, err := service.CategoryCache.GetCategoryBatch(ctx)

	if err == nil {
		categories = append(categories, category)

		service.CategoryCache.SetCategoryBatch(ctx, categories)
	}

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

	service.CategoryCache.Delete(context.Background(), category)

	return helper.ToCategoryResponse(category), nil
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) error {

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return errors.New("category not found")
	}

	service.CategoryRepository.Delete(ctx, category)

	service.CategoryCache.Delete(context.Background(), category)

	return nil
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) (*web.CategoryResponse, error) {

	categoryCache, err := service.CategoryCache.GetCategory(context.Background(), categoryId)

	log.Println("category cache ", categoryCache)

	if err == nil {
		return helper.ToCategoryResponse(categoryCache), nil
	}

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return nil, err
	}

	service.CategoryCache.SetCategory(context.Background(), category)

	return helper.ToCategoryResponse(category), nil
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []*web.CategoryResponse {

	categoriesCache, err := service.CategoryCache.GetCategoryBatch(context.Background())

	if err == nil {
		return helper.ToCategoryResponses(categoriesCache)
	}

	categories := service.CategoryRepository.FindAll(ctx)

	if len(categories) > 0 {
		service.CategoryCache.SetCategoryBatch(context.Background(), categories)
	}

	return helper.ToCategoryResponses(categories)
}
