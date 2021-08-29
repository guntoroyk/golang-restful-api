package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"github.com/guntoroyk/golang-restful-api/model/web"
	"github.com/guntoroyk/golang-restful-api/repository/mocks"
)

func TestCategoryServiceImpl_FindAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryRepo.EXPECT().FindAll(context.Background()).Return([]domain.Category{{
		Id:   1,
		Name: "Computer",
	}, {
		Id:   2,
		Name: "Gadget",
	}})

	categoryService := NewCategoryService(mockCategoryRepo)

	categories := categoryService.FindAll(context.Background())

	fmt.Println("categories", categories)
	got := len(categories)
	want := 2
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCategoryServiceImpl_FindById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	category := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	mockCategoryRepo.EXPECT().FindById(context.Background(), category.Id).Return(domain.Category{Id: 1, Name: "Computer"}, nil)

	categoryService := NewCategoryService(mockCategoryRepo)

	categoryResult := categoryService.FindById(context.Background(), 1)

	fmt.Println("category", categoryResult)
	got := categoryResult
	want := web.CategoryResponse{Id: 1, Name: "Computer"}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCategoryServiceImpl_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	category := domain.Category{
		Name: "Computer",
	}
	mockCategoryRepo.EXPECT().Save(context.Background(), category).Return(domain.Category{Id: 1, Name: "Computer"})
	categoryService := NewCategoryService(mockCategoryRepo)

	request := web.CategoryCreateRequest{Name: "Computer"}

	categoryResult := categoryService.Create(context.Background(), request)
	fmt.Println("CategoryResult", categoryResult)

	got := categoryResult
	want := web.CategoryResponse{Id: 1, Name: "Computer"}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCategoryServiceImpl_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	category := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	mockCategoryRepo.EXPECT().FindById(context.Background(), category.Id).Return(domain.Category{Id: 1, Name: "Computer"}, nil)
	mockCategoryRepo.EXPECT().Delete(context.Background(), category)

	categoryService := NewCategoryService(mockCategoryRepo)

	categoryService.Delete(context.Background(), 1)

	//fmt.Println("category", categoryResult)
	//got := categoryResult
	//want := web.CategoryResponse{Id: 1, Name: "Computer"}
	//
	//if got != want {
	//	t.Errorf("got %q want %q", got, want )
	//}
}

func TestCategoryServiceImpl_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	oldCategory := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	mockCategoryRepo.EXPECT().FindById(context.Background(), 1).Return(oldCategory, nil)

	newCategory := domain.Category{
		Id:   1,
		Name: "Gadget",
	}

	mockCategoryRepo.EXPECT().Update(context.Background(), newCategory).Return(newCategory)

	categoryService := NewCategoryService(mockCategoryRepo)

	request := web.CategoryUpdateRequest{Id: 1, Name: "Gadget"}

	categoryResult := categoryService.Update(context.Background(), request)
	fmt.Println("CategoryResult", categoryResult)

	got := categoryResult
	want := web.CategoryResponse{Id: 1, Name: "Gadget"}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
