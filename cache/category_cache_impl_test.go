package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"reflect"
	"testing"
	"time"
)

var (
	jsonMarshal = json.Marshal
)

func fakemarshal(v interface{}) ([]byte, error) {
	return []byte{}, errors.New("marshalling failed")
}

func restoremarshal(replace func(v interface{}) ([]byte, error)) {
	jsonMarshal = replace
}

func TestCategoryCacheImpl_GetCategory(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		ctx        context.Context
		categoryId int
	}

	db, mock := redismock.NewClientMock()
	defer mock.ExpectationsWereMet()

	categoryId := 1
	key := fmt.Sprintf("category:%d", categoryId)

	var emptyCategory domain.Category
	categoryFound := domain.Category{
		Id:   1,
		Name: "Computer",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Category
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "category cache is nil",
			fields: func() fields {
				mock.ExpectGet(key).RedisNil()

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background(), categoryId: categoryId},
			want:    emptyCategory,
			wantErr: true,
		},
		{
			name: "category cache is found",
			fields: func() fields {
				mock.ExpectGet(key).SetVal("{\"Id\":1,\"Name\":\"Computer\"}")

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background(), categoryId: categoryId},
			want:    categoryFound,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CategoryCacheImpl{
				redisClient: tt.fields.redisClient,
			}
			got, err := c.GetCategory(tt.args.ctx, tt.args.categoryId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCategory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryCacheImpl_GetCategoryBatch(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		ctx context.Context
	}

	key := "categories"

	db, mock := redismock.NewClientMock()

	defer mock.ExpectationsWereMet()

	var emptyCategories []domain.Category
	category := domain.Category{Id: 1, Name: "Computer"}

	categories := []domain.Category{category}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Category
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "categories cache is nil",
			fields: func() fields {
				mock.ExpectGet(key).RedisNil()

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background()},
			want:    emptyCategories,
			wantErr: true,
		},
		{
			name: "categories cache is found",
			fields: func() fields {
				mock.ExpectGet(key).SetVal("[{\"Id\":1,\"Name\":\"Computer\"}]")

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background()},
			want:    categories,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CategoryCacheImpl{
				redisClient: tt.fields.redisClient,
			}
			got, err := c.GetCategoryBatch(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategoryBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCategoryBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryCacheImpl_Delete(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		ctx      context.Context
		category domain.Category
	}

	db, mock := redismock.NewClientMock()
	//defer mock.ExpectationsWereMet()

	category := domain.Category{Id: 1, Name: "Computer"}

	key := fmt.Sprintf("category:%d", category.Id)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete category cache error",
			fields: func() fields {
				mock.ExpectDel(key).SetErr(errors.New("FAIL"))

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background(), category: category},
			wantErr: true,
		}, {
			name: "delete category cache success",
			fields: func() fields {
				mock.ExpectDel(key).SetVal(1)

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background(), category: category},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CategoryCacheImpl{
				redisClient: tt.fields.redisClient,
			}
			if err := c.Delete(tt.args.ctx, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryCacheImpl_SetCategory(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		ctx      context.Context
		category domain.Category
	}

	db, mock := redismock.NewClientMock()
	defer mock.ExpectationsWereMet()

	//storedMarshal := jsonMarshal
	//jsonMarshal = fakemarshal
	//defer restoremarshal(storedMarshal)

	category := domain.Category{Id: 1, Name: "Computer"}
	json, _ := json.Marshal(category)

	key := fmt.Sprintf("category:%d", category.Id)
	expiration := 5 * time.Minute

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create category success",
			fields: func() fields {
				mock.ExpectSet(key, json, expiration).SetVal("OK")

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background(), category: category},
			wantErr: false,
		},
		{
			name: "create category error",
			fields: func() fields {
				mock.ExpectSet(key, json, expiration).SetErr(errors.New("FAIL"))

				return fields{redisClient: db}
			}(),
			args:    args{ctx: context.Background(), category: category},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CategoryCacheImpl{
				redisClient: tt.fields.redisClient,
			}
			if err := c.SetCategory(tt.args.ctx, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("SetCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
