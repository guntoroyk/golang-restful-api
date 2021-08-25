package main

import (
	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"golang-restful-api/controller"
	"golang-restful-api/db"
	"golang-restful-api/exception"
	"golang-restful-api/helper"
	"golang-restful-api/repository"
	"golang-restful-api/service"
	"net/http"
)

func main() {

	dbConfig := db.Config{
		User: "postgres",
		Password: "admin",
		DBName: "golang_restful_api",
		Port: 5433,
		Host: "localhost",
		SSLMode: "disable",
	}

	db := db.NewDB(dbConfig)
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := httprouter.New()

	router.GET("/api/categories", categoryController.FindAll)
	router.GET("/api/categories/:categoryId", categoryController.FindById)
	router.POST("/api/categories", categoryController.Create)
	router.PUT("/api/categories/:categoryId", categoryController.Update)
	router.DELETE("/api/categories/:categoryId", categoryController.Delete)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr: "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
