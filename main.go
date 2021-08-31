package main

import (
	"github.com/guntoroyk/golang-restful-api/cache"
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
		User:     "postgres",
		Password: "admin",
		DBName:   "golang_restful_api",
		Port:     5433,
		Host:     "localhost",
		SSLMode:  "disable",
	}

	redisConfig := app.RedisConfig{
		Addr: "localhost:6379",
		DB:   0,
	}

	db := app.NewDB(dbConfig)
	redisClient := app.NewRedisClient(redisConfig)
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository(db, validate)
	categoryCache := cache.NewCategoryCache(redisClient)
	categoryService := service.NewCategoryService(categoryRepository, categoryCache)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
