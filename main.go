package main

import (
	"github.com/guntoroyk/golang-restful-api/cache"
	"github.com/guntoroyk/golang-restful-api/messaging"
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
	db := app.NewDB(dbConfig)

	redisConfig := app.RedisConfig{
		Addr: "localhost:6379",
		DB:   0,
	}
	redisClient := app.NewRedisClient(redisConfig)

	producerConfig := messaging.ProducerConfig{
		NsqdAddress: "nsqd:4150",
	}
	messageProducer := messaging.NewProducer(producerConfig)

	validate := validator.New()

	categoryRepository := repository.NewCategoryRepository(db, validate)
	categoryCache := cache.NewCategoryCache(redisClient)
	categoryService := service.NewCategoryService(categoryRepository, categoryCache, messageProducer)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
