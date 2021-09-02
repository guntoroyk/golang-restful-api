package service

import (
	"context"
	"errors"
	"github.com/guntoroyk/golang-restful-api/cache"
	"github.com/guntoroyk/golang-restful-api/helper"
	"github.com/guntoroyk/golang-restful-api/messaging"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"github.com/guntoroyk/golang-restful-api/model/web"
	"github.com/guntoroyk/golang-restful-api/repository"
	"log"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	CategoryCache      cache.CategoryCache
	Producer           messaging.Producer
}

func NewCategoryService(categoryRepository repository.CategoryRepository, categoryCache cache.CategoryCache, p messaging.Producer) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		CategoryCache:      categoryCache,
		Producer:           p,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) *web.CategoryResponse {
	var resp *web.CategoryResponse

	defer func() {
		message := web.CategoryProducerMessage{
			Event:          "view",
			CategoryDetail: resp,
		}

		if err := service.Producer.Publish("category_view", message); err != nil {
			log.Println("failed to publish message data, err: ", err.Error())
		}
	}()

	category := domain.Category{
		Name: request.Name,
	}

	category = service.CategoryRepository.Save(ctx, category)

	categories, err := service.CategoryCache.GetCategoryBatch(ctx)

	if err == nil {
		categories = append(categories, category)

		service.CategoryCache.SetCategoryBatch(ctx, categories)
	}

	resp = helper.ToCategoryResponse(category)

	return resp
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) (*web.CategoryResponse, error) {

	category, err := service.CategoryRepository.FindById(ctx, request.Id)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return nil, errors.New("category not found")
	}

	category.Name = request.Name

	category = service.CategoryRepository.Update(ctx, category)

	service.CategoryCache.Delete(ctx, category)

	return helper.ToCategoryResponse(category), nil
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) error {

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return errors.New("category not found")
	}

	service.CategoryRepository.Delete(ctx, category)

	service.CategoryCache.Delete(ctx, category)

	return nil
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) (*web.CategoryResponse, error) {
	var resp *web.CategoryResponse

	defer func() {
		message := web.CategoryProducerMessage{
			Event:          "view",
			CategoryDetail: resp,
		}

		if err := service.Producer.Publish("category_view", message); err != nil {
			log.Println("failed to publish message data, err: ", err.Error())
		}
	}()

	categoryCache, err := service.CategoryCache.GetCategory(context.Background(), categoryId)

	log.Println("category cache ", categoryCache)

	if err == nil {
		resp = helper.ToCategoryResponse(categoryCache)
		return resp, nil
	}

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		//panic(exception.NewNotFoundError(err.Error()))
		return nil, err
	}

	service.CategoryCache.SetCategory(ctx, category)

	resp = helper.ToCategoryResponse(category)

	return resp, nil
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []*web.CategoryResponse {

	categoriesCache, err := service.CategoryCache.GetCategoryBatch(ctx)

	if err == nil {
		return helper.ToCategoryResponses(categoriesCache)
	}

	categories := service.CategoryRepository.FindAll(ctx)

	if len(categories) > 0 {
		service.CategoryCache.SetCategoryBatch(ctx, categories)
	}

	return helper.ToCategoryResponses(categories)
}
