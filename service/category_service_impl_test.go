package service

import (
	"context"
	"fmt"
	"github.com/guntoroyk/golang-restful-api/repository"
	"reflect"
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

func TestCategoryServiceImpl_Create1(t *testing.T) {
	type fields struct {
		CategoryRepository *mocks.MockCategoryRepository
	}
	type args struct {
		ctx     context.Context
		request web.CategoryCreateRequest
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	mockCategoryRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(domain.Category{Id: 1, Name: "Computer"})

	tests := []struct {
		name   string
		fields fields
		args   args
		want   web.CategoryResponse
	}{
		// TODO: Add test cases.
		{
			name: "Create category with name Computer",
			fields: fields{
				CategoryRepository: mockCategoryRepo,
			},
			args: args{ctx: context.Background(), request: web.CategoryCreateRequest{Name: "Computer"}},
			want: web.CategoryResponse{
				Id:   1,
				Name: "Computer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
			}
			if got := service.Create(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryServiceImpl_Update1(t *testing.T) {
	type fields struct {
		CategoryRepository *mocks.MockCategoryRepository
	}
	type args struct {
		ctx     context.Context
		request web.CategoryUpdateRequest
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	oldCategory := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	mockCategoryRepo.EXPECT().FindById(gomock.Any(), 1).Return(oldCategory, nil)

	updatedCategory := domain.Category{
		Id:   1,
		Name: "Gadget",
	}

	mockCategoryRepo.EXPECT().Update(gomock.Any(), updatedCategory).Return(updatedCategory)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   web.CategoryResponse
	}{
		// TODO: Add test cases.
		{
			name: "Update category",
			fields: fields{
				CategoryRepository: mockCategoryRepo,
			},
			args: args{ctx: context.Background(), request: web.CategoryUpdateRequest{Id: 1, Name: "Gadget"}},
			want: web.CategoryResponse{Id: updatedCategory.Id, Name: updatedCategory.Name},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
			}
			if got := service.Update(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryServiceImpl_Delete1(t *testing.T) {
	type fields struct {
		CategoryRepository *mocks.MockCategoryRepository
	}
	type args struct {
		ctx        context.Context
		categoryId int
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	categoryToDelete := domain.Category{Id: 1, Name: "Computer"}

	mockCategoryRepo.EXPECT().FindById(gomock.Any(), categoryToDelete.Id).Return(categoryToDelete, nil)
	mockCategoryRepo.EXPECT().Delete(gomock.Any(), categoryToDelete)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name:   "Delete category",
			fields: fields{CategoryRepository: mockCategoryRepo},
			args:   args{ctx: context.Background(), categoryId: categoryToDelete.Id},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
			}

			service.Delete(tt.args.ctx, tt.args.categoryId)
		})
	}
}

func TestCategoryServiceImpl_FindById1(t *testing.T) {
	type fields struct {
		CategoryRepository repository.CategoryRepository
	}
	type args struct {
		ctx        context.Context
		categoryId int
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	category := domain.Category{Id: 1, Name: "Computer"}

	mockCategoryRepo.EXPECT().FindById(gomock.Any(), category.Id).Return(category, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   web.CategoryResponse
	}{
		// TODO: Add test cases.
		{
			name:   "Find a category by Id",
			fields: fields{CategoryRepository: mockCategoryRepo},
			args:   args{ctx: context.Background(), categoryId: category.Id},
			want:   web.CategoryResponse{Id: category.Id, Name: category.Name},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
			}
			if got := service.FindById(tt.args.ctx, tt.args.categoryId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryServiceImpl_FindAll1(t *testing.T) {
	type fields struct {
		CategoryRepository *mocks.MockCategoryRepository
	}
	type args struct {
		ctx context.Context
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)

	category := domain.Category{Id: 1, Name: "Computer"}
	categories := []domain.Category{category}

	mockCategoryRepo.EXPECT().FindAll(gomock.Any()).Return(categories)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []web.CategoryResponse
	}{
		// TODO: Add test cases.
		{
			name:   "Find all categories",
			fields: fields{CategoryRepository: mockCategoryRepo},
			args:   args{ctx: context.Background()},
			want: []web.CategoryResponse{{
				Id:   category.Id,
				Name: category.Name,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
			}
			if got := service.FindAll(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
