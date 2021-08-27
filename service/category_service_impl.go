package service

import (
	"context"
	"github.com/guntoroyk/golang-restful-api/exception"
	"github.com/guntoroyk/golang-restful-api/helper"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"github.com/guntoroyk/golang-restful-api/model/web"
	"github.com/guntoroyk/golang-restful-api/repository"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	//DB *sql.DB
	//Validate *validator.Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		//DB: DB,
		//Validate: validate,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	//err := service.Validate.Struct(request)
	//helper.PanicIfError(err)
	//
	//tx, err := service.DB.Begin()
	//helper.PanicIfError(err)
	//defer helper.CommitOrRollback(tx)

	category := domain.Category{
		Name: request.Name,
	}

	category = service.CategoryRepository.Save(ctx, category)

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	//err := service.Validate.Struct(request)
	//helper.PanicIfError(err)

	//tx, err := service.DB.Begin()
	//if err != nil {
	//	panic(exception.NewNotFoundError(err.Error()))
	//}
	//
	//defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	category.Name = request.Name

	category = service.CategoryRepository.Update(ctx, category)

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	//tx, err := service.DB.Begin()
	//helper.PanicIfError(err)
	//defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.CategoryRepository.Delete(ctx, category)
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	//tx, err := service.DB.Begin()
	//helper.PanicIfError(err)
	//defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	//tx, err := service.DB.Begin()
	//helper.PanicIfError(err)
	//defer helper.CommitOrRollback(tx)

	categories := service.CategoryRepository.FindAll(ctx)

	return helper.ToCategoryResponses(categories)
}
