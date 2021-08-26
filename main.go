package main

import (
	"github.com/go-playground/validator"
	_ "github.com/lib/pq"
	"golang-restful-api/app"
	"golang-restful-api/controller"
	"golang-restful-api/helper"
	"golang-restful-api/middleware"
	"golang-restful-api/repository"
	"golang-restful-api/service"
	"net/http"
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
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr: "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
