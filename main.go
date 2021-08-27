package main

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/guntoroyk/golang-restful-api/app"
	"github.com/guntoroyk/golang-restful-api/controller"
	"github.com/guntoroyk/golang-restful-api/helper"
	"github.com/guntoroyk/golang-restful-api/middleware"
	"github.com/guntoroyk/golang-restful-api/repository"
	"github.com/guntoroyk/golang-restful-api/service"
	_ "github.com/lib/pq"
)

func main() {

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
	categoryService := service.NewCategoryService(categoryRepository)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr: "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
