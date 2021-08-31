package main

import (
	"context"
	"fmt"
	"github.com/guntoroyk/golang-restful-api/model/web"
	"github.com/guntoroyk/golang-restful-api/service"
	"log"
	"net"

	"github.com/go-playground/validator"
	"github.com/guntoroyk/golang-restful-api/app"
	"github.com/guntoroyk/golang-restful-api/grpc/proto"
	"github.com/guntoroyk/golang-restful-api/repository"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type server struct {
	categoryService service.CategoryService
}

func (s server) CreateCategory(ctx context.Context, request *proto.CreateCategoryRequest) (*proto.CategoryResponse, error) {
	categoryCreated := s.categoryService.Create(ctx, web.CategoryCreateRequest{Name: request.Name})

	res := &proto.CategoryResponse{Category: &proto.Category{Id: int32(categoryCreated.Id), Name: categoryCreated.Name}}

	return res, nil
}

func (s server) UpdateCategory(ctx context.Context, request *proto.UpdateCategoryRequest) (*proto.CategoryResponse, error) {
	categoryUpdated, err := s.categoryService.Update(ctx, web.CategoryUpdateRequest{Id: int(request.Category.Id), Name: request.Category.Name})

	if err != nil {
		return nil, err
	}
	res := &proto.CategoryResponse{Category: &proto.Category{Id: int32(categoryUpdated.Id), Name: categoryUpdated.Name}}

	return res, nil
}

func (s server) GetCategory(ctx context.Context, request *proto.GetCategoryRequest) (*proto.CategoryResponse, error) {
	//categoryFound, err := s.categoryRepository.FindById(ctx, int(request.CategoryId))
	//if err != nil {
	//	return nil, err
	//}

	categoryFound, err := s.categoryService.FindById(ctx, int(request.CategoryId))

	if err != nil {
		return nil, err
	}

	res := &proto.CategoryResponse{Category: &proto.Category{Id: int32(categoryFound.Id), Name: categoryFound.Name}}

	return res, nil
}

func (s server) GetAllCategory(ctx context.Context, request *proto.GetAllCategoryRequest) (*proto.GetAllCategoryResponse, error) {
	fmt.Printf("GetAllCategory function was invoked with %v", request)

	categoriesRes := s.categoryService.FindAll(ctx)

	fmt.Println("CategoriesREs", categoriesRes)
	var categories []*proto.Category

	for _, category := range categoriesRes {
		categories = append(categories, &proto.Category{
			Id:   int32(category.Id),
			Name: category.Name,
		})
	}

	return &proto.GetAllCategoryResponse{
		Categories: categories,
	}, nil
}

func (s server) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*proto.DeleteCategoryResponse, error) {
	err := s.categoryService.Delete(ctx, int(request.CategoryId))

	if err != nil {
		return nil, err
	}

	res := &proto.DeleteCategoryResponse{Message: "Category successfully deleted"}

	return res, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("GRPC Server starting...")

	dbConfig := app.Config{
		User:     "postgres",
		Password: "admin",
		DBName:   "golang_restful_api",
		Port:     5433,
		Host:     "localhost",
		SSLMode:  "disable",
	}

	db := app.NewDB(dbConfig)
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository(db, validate)
	categoryService := service.NewCategoryService(categoryRepository)

	list, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterCategoryServiceServer(s, &server{
		categoryService: categoryService,
	})

	if err := s.Serve(list); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
