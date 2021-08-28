package main

import (
	"context"
	"fmt"
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
	categoryRepository repository.CategoryRepository
}

func (s server) CreateCategory(ctx context.Context, request *proto.CreateCategoryRequest) (*proto.CreateCategoryResponse, error) {
	panic("implement me")
}

func (s server) UpdateCategory(ctx context.Context, request *proto.UpdateCategoryRequest) (*proto.UpdateCategoryResponse, error) {
	panic("implement me")
}

func (s server) GetCategory(ctx context.Context, request *proto.GetCategoryRequest) (*proto.GetCategoryResponse, error) {
	categoryFound, err := s.categoryRepository.FindById(ctx, int(request.CategoryId))
	if err != nil {
		return nil, err
	}

	res := &proto.GetCategoryResponse{Category: &proto.Category{Id: int32(categoryFound.Id), Name: categoryFound.Name}}

	return res, nil
}

func (s server) GetAllCategory(ctx context.Context, request *proto.GetAllCategoryRequest) (*proto.GetAllCategoryResponse, error) {
	fmt.Printf("GetAllCategory function was invoked with %v", request)

	categoriesRes := s.categoryRepository.FindAll(ctx)

	fmt.Println("CategoriesREs", categoriesRes)
	var categories []*proto.Category

	for _, category := range categoriesRes {
		categories = append(categories, &proto.Category{
			Id: int32(category.Id),
			Name: category.Name,
		})
	}

	return &proto.GetAllCategoryResponse{
		Categories: categories,
	}, nil
}

func (s server) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*proto.DeleteCategoryResponse, error) {
	panic("implement me")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("GRPC Server starting...")

	dbConfig := app.Config{
		User: "postgres",
		Password: "admin",
		DBName: "golang_restful_api",
		Port: 5433,
		Host: "localhost",
		SSLMode: "disable",
	}

	db := app.NewDB(dbConfig)
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository(db, validate)
	//categoryService := service.NewCategoryService(categoryRepository)


	list, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterCategoryServiceServer(s, &server{
		categoryRepository: categoryRepository,
	})

	if err := s.Serve(list); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
