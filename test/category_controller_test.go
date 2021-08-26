package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang-restful-api/app"
	"golang-restful-api/controller"
	"golang-restful-api/middleware"
	"golang-restful-api/model/domain"
	"golang-restful-api/repository"
	"golang-restful-api/service"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func setupDBTest() *sql.DB {
	dbConfig := app.Config{
		User: "postgres",
		Password: "admin",
		DBName: "golang_restful_api_test",
		Port: 5433,
		Host: "localhost",
		SSLMode: "disable",
	}

	db := app.NewDB(dbConfig)

	return db
}

func setupRouter(db *sql.DB) http.Handler {

	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB)  {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupDBTest()
	router := setupRouter(db)
	truncateCategory(db)

	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
	//fmt.Println("response.Body", &response.Body)
	//
	//body, _ := io.ReadAll(response.Body)
	//var responseBody map[string]interface{}
	//json.Unmarshal(body, responseBody)
	//
	//fmt.Println("responBody", responseBody)
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupDBTest()
	router := setupRouter(db)
	truncateCategory(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)
}


func TestUpdateCategorySuccess(t *testing.T) {
	db := setupDBTest()
	router := setupRouter(db)
	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	fmt.Println("categoryId", category.Id)

	requestBody := strings.NewReader(`{"name": "Computer"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories" + strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-Key", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)
}

func TestUpdateCategoryFailed(t *testing.T) {

}


