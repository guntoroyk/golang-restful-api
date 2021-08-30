package repository

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator"
	"github.com/guntoroyk/golang-restful-api/model/domain"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCategoryRepositoryImpl_Save(t *testing.T) {
	type fields struct {
		DB       *sql.DB
		Validate *validator.Validate
	}
	type args struct {
		ctx      context.Context
		category domain.Category
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	defer mock.ExpectationsWereMet()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery("insert into category(name) values ($1) returning id").WithArgs("Computer").WillReturnRows(rows)
	mock.ExpectCommit()

	validate := validator.New()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   domain.Category
	}{
		// TODO: Add test cases.
		{
			name: "Create new category",
			fields: fields{
				DB:       db,
				Validate: validate,
			},
			args: args{
				ctx: context.Background(),
				category: domain.Category{
					Name: "Computer",
				},
			},
			want: domain.Category{
				Id:   1,
				Name: "Computer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &CategoryRepositoryImpl{
				DB:       tt.fields.DB,
				Validate: tt.fields.Validate,
			}
			if got := repository.Save(tt.args.ctx, tt.args.category); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryRepositoryImpl_FindById(t *testing.T) {
	type fields struct {
		DB       *sql.DB
		Validate *validator.Validate
	}
	type args struct {
		ctx        context.Context
		categoryId int
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Computer")

	mock.ExpectBegin()
	mock.ExpectQuery("select id, name from category where id = $1").WithArgs(1).WillReturnRows(rows)
	mock.ExpectCommit()

	validate := validator.New()

	category := domain.Category{Id: 1, Name: "Computer"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Category
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Get category by ID",
			fields:  fields{DB: db, Validate: validate},
			args:    args{ctx: context.Background(), categoryId: category.Id},
			want:    category,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &CategoryRepositoryImpl{
				DB:       tt.fields.DB,
				Validate: tt.fields.Validate,
			}
			got, err := repository.FindById(tt.args.ctx, tt.args.categoryId)
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
