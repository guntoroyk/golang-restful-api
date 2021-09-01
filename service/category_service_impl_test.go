package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/guntoroyk/golang-restful-api/cache"
	cache_mocks "github.com/guntoroyk/golang-restful-api/cache/mocks"
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

	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)
	mockCategoryCache.EXPECT().GetCategoryBatch(context.Background()).Return([]domain.Category{{
		Id:   1,
		Name: "Computer",
	}, {
		Id:   2,
		Name: "Gadget",
	}}, nil)

	categoryService := NewCategoryService(mockCategoryRepo, mockCategoryCache)

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

	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)
	mockCategoryCache.EXPECT().GetCategory(context.Background(), category.Id).Return(domain.Category{Id: 1, Name: "Computer"}, nil)

	categoryService := NewCategoryService(mockCategoryRepo, mockCategoryCache)

	categoryResult, _ := categoryService.FindById(context.Background(), 1)

	fmt.Println("category", categoryResult)
	got := categoryResult
	want := &web.CategoryResponse{Id: 1, Name: "Computer"}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCategoryServiceImpl_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	category := domain.Category{
		Name: "Computer",
	}
	mockCategoryRepo.EXPECT().Save(context.Background(), category).Return(domain.Category{Id: 1, Name: "Computer"})
	mockCategoryCache.EXPECT().SetCategory(context.Background(), category).Return(domain.Category{Id: 1, Name: "Computer"})

	categoryService := NewCategoryService(mockCategoryRepo, mockCategoryCache)

	request := web.CategoryCreateRequest{Name: "Computer"}

	categoryResult := categoryService.Create(context.Background(), request)
	fmt.Println("CategoryResult", categoryResult)

	got := categoryResult
	want := &web.CategoryResponse{Id: 1, Name: "Computer"}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCategoryServiceImpl_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	category := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	mockCategoryRepo.EXPECT().FindById(context.Background(), category.Id).Return(domain.Category{Id: 1, Name: "Computer"}, nil)
	mockCategoryRepo.EXPECT().Delete(context.Background(), category)
	mockCategoryCache.EXPECT().GetCategory(context.Background(), category.Id).Return(domain.Category{Id: 1, Name: "Computer"}, nil)
	mockCategoryCache.EXPECT().Delete(context.Background(), category)

	categoryService := NewCategoryService(mockCategoryRepo, mockCategoryCache)

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
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

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
	mockCategoryCache.EXPECT().Delete(context.Background(), oldCategory).Return(newCategory)

	categoryService := NewCategoryService(mockCategoryRepo, mockCategoryCache)

	request := web.CategoryUpdateRequest{Id: 1, Name: "Gadget"}

	categoryResult, _ := categoryService.Update(context.Background(), request)
	fmt.Println("CategoryResult", categoryResult)

	got := categoryResult
	want := &web.CategoryResponse{Id: 1, Name: "Gadget"}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCategoryServiceImpl_FindAll1(t *testing.T) {
	type fields struct {
		CategoryRepository repository.CategoryRepository
		CategoryCache      cache.CategoryCache
	}
	type args struct {
		ctx context.Context
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	category := domain.Category{Id: 1, Name: "Computer"}
	categories := []domain.Category{category}
	var emptyCategories []domain.Category

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*web.CategoryResponse
	}{
		// test cases.
		{
			name: "Find all categories from cache",
			fields: func() fields {
				mockCategoryCache.EXPECT().GetCategoryBatch(gomock.Any()).Times(1).Return(categories, nil)

				mockCategoryRepo.EXPECT().FindAll(gomock.Any()).Times(0)

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args: args{ctx: context.Background()},
			want: []*web.CategoryResponse{{
				Id:   category.Id,
				Name: category.Name,
			}},
		},
		{
			name: "Find all categories from database",
			fields: func() fields {
				mockCategoryCache.EXPECT().GetCategoryBatch(gomock.Any()).Times(1).Return(emptyCategories, errors.New("redis: nil"))

				mockCategoryRepo.EXPECT().FindAll(gomock.Any()).Times(1).Return(categories)

				mockCategoryCache.EXPECT().SetCategoryBatch(gomock.Any(), categories).Times(1).Return(nil)

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args: args{ctx: context.Background()},
			want: []*web.CategoryResponse{{
				Id:   category.Id,
				Name: category.Name,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
				CategoryCache:      tt.fields.CategoryCache,
			}
			if got := service.FindAll(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryServiceImpl_FindById1(t *testing.T) {
	type fields struct {
		CategoryRepository repository.CategoryRepository
		CategoryCache      cache.CategoryCache
	}
	type args struct {
		ctx        context.Context
		categoryId int
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	category := domain.Category{Id: 1, Name: "Computer"}
	var emptyCategory domain.Category

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *web.CategoryResponse
		wantErr bool
	}{
		// test cases.
		{
			name: "FindById from cache",
			fields: func() fields {
				mockCategoryCache.EXPECT().GetCategory(gomock.Any(), category.Id).Times(1).Return(category, nil)

				mockCategoryRepo.EXPECT().FindById(gomock.Any(), category.Id).Times(0)

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args:    args{ctx: context.Background(), categoryId: category.Id},
			want:    &web.CategoryResponse{Id: category.Id, Name: category.Name},
			wantErr: false,
		}, {
			name: "FindById from db",
			fields: func() fields {
				mockCategoryCache.EXPECT().GetCategory(gomock.Any(), category.Id).Times(1).Return(emptyCategory, errors.New("redis: nil"))

				mockCategoryRepo.EXPECT().FindById(gomock.Any(), category.Id).Times(1).Return(category, nil)

				mockCategoryCache.EXPECT().SetCategory(gomock.Any(), category).Times(1).Return(nil)

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args:    args{ctx: context.Background(), categoryId: category.Id},
			want:    &web.CategoryResponse{Id: category.Id, Name: category.Name},
			wantErr: false,
		}, {
			name: "FindById error not found",
			fields: func() fields {
				mockCategoryCache.EXPECT().GetCategory(gomock.Any(), category.Id).Times(1).Return(emptyCategory, errors.New("redis: nil"))

				mockCategoryRepo.EXPECT().FindById(gomock.Any(), category.Id).Times(1).Return(emptyCategory, errors.New("category not found"))

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args:    args{ctx: context.Background(), categoryId: category.Id},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
				CategoryCache:      tt.fields.CategoryCache,
			}
			got, err := service.FindById(tt.args.ctx, tt.args.categoryId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryServiceImpl_Delete1(t *testing.T) {
	type fields struct {
		CategoryRepository repository.CategoryRepository
		CategoryCache      cache.CategoryCache
	}
	type args struct {
		ctx        context.Context
		categoryId int
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	categoryToDelete := domain.Category{Id: 1, Name: "Computer"}
	var emptyCategory domain.Category

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// test cases.
		{
			name: "delete category success",
			fields: func() fields {
				mockCategoryRepo.EXPECT().FindById(gomock.Any(), categoryToDelete.Id).Return(categoryToDelete, nil)
				mockCategoryRepo.EXPECT().Delete(gomock.Any(), categoryToDelete)
				mockCategoryCache.EXPECT().Delete(gomock.Any(), categoryToDelete).Return(nil)

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args: args{
				ctx:        context.Background(),
				categoryId: categoryToDelete.Id,
			},
			wantErr: false,
		}, {
			name: "delete category failed",
			fields: func() fields {
				mockCategoryRepo.EXPECT().FindById(gomock.Any(), categoryToDelete.Id).Return(emptyCategory, errors.New("category not found"))

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args: args{
				ctx:        context.Background(),
				categoryId: categoryToDelete.Id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
				CategoryCache:      tt.fields.CategoryCache,
			}
			if err := service.Delete(tt.args.ctx, tt.args.categoryId); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryServiceImpl_Update2(t *testing.T) {
	type fields struct {
		CategoryRepository repository.CategoryRepository
		CategoryCache      cache.CategoryCache
	}
	type args struct {
		ctx     context.Context
		request web.CategoryUpdateRequest
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	oldCategory := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	updatedCategory := domain.Category{
		Id:   1,
		Name: "Gadget",
	}

	var emptyCategory domain.Category

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *web.CategoryResponse
		wantErr bool
	}{
		// test cases.
		{
			name: "Update category success",
			fields: func() fields {
				mockCategoryRepo.EXPECT().FindById(gomock.Any(), oldCategory.Id).Times(1).Return(oldCategory, nil)

				mockCategoryRepo.EXPECT().Update(gomock.Any(), updatedCategory).Times(1).Return(updatedCategory)

				mockCategoryCache.EXPECT().Delete(gomock.Any(), updatedCategory)

				return fields{
					CategoryRepository: mockCategoryRepo,
					CategoryCache:      mockCategoryCache,
				}
			}(),
			args:    args{ctx: context.Background(), request: web.CategoryUpdateRequest{Id: updatedCategory.Id, Name: updatedCategory.Name}},
			want:    &web.CategoryResponse{Id: updatedCategory.Id, Name: updatedCategory.Name},
			wantErr: false,
		},
		{
			name: "Update category failed",
			fields: func() fields {
				mockCategoryRepo.EXPECT().FindById(gomock.Any(), oldCategory.Id).Times(1).Return(emptyCategory, errors.New("category not found"))

				return fields{
					CategoryRepository: mockCategoryRepo,
					CategoryCache:      mockCategoryCache,
				}
			}(),
			args:    args{ctx: context.Background(), request: web.CategoryUpdateRequest{Id: updatedCategory.Id, Name: updatedCategory.Name}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
				CategoryCache:      tt.fields.CategoryCache,
			}
			got, err := service.Update(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryServiceImpl_Create2(t *testing.T) {
	type fields struct {
		CategoryRepository repository.CategoryRepository
		CategoryCache      cache.CategoryCache
	}
	type args struct {
		ctx     context.Context
		request web.CategoryCreateRequest
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCategoryRepo := mocks.NewMockCategoryRepository(mockCtrl)
	mockCategoryCache := cache_mocks.NewMockCategoryCache(mockCtrl)

	newCategory := domain.Category{Id: 2, Name: "Computer"}

	categories := []domain.Category{
		{
			Id:   1,
			Name: "Gadget",
		},
	}

	var emptyCategories []domain.Category

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *web.CategoryResponse
	}{
		// TODO: Add test cases.
		{
			name: "Create newCategory success",
			fields: func() fields {
				mockCategoryRepo.EXPECT().Save(gomock.Any(), domain.Category{Name: newCategory.Name}).Times(1).Return(newCategory)
				mockCategoryCache.EXPECT().GetCategoryBatch(gomock.Any()).Times(1).Return(emptyCategories, errors.New("redis: nil"))

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args: args{ctx: context.Background(), request: web.CategoryCreateRequest{Name: newCategory.Name}},
			want: &web.CategoryResponse{
				Id:   newCategory.Id,
				Name: newCategory.Name,
			},
		},
		{
			name: "Create newCategory success and update existing cache",
			fields: func() fields {
				mockCategoryRepo.EXPECT().Save(gomock.Any(), domain.Category{Name: newCategory.Name}).Times(1).Return(newCategory)
				mockCategoryCache.EXPECT().GetCategoryBatch(gomock.Any()).Times(1).Return(categories, nil)

				categories = append(categories, newCategory)

				mockCategoryCache.EXPECT().SetCategoryBatch(gomock.Any(), categories)

				return fields{CategoryRepository: mockCategoryRepo, CategoryCache: mockCategoryCache}
			}(),
			args: args{ctx: context.Background(), request: web.CategoryCreateRequest{Name: newCategory.Name}},
			want: &web.CategoryResponse{
				Id:   newCategory.Id,
				Name: newCategory.Name,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &CategoryServiceImpl{
				CategoryRepository: tt.fields.CategoryRepository,
				CategoryCache:      tt.fields.CategoryCache,
			}
			if got := service.Create(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
