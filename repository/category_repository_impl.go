package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator"
	"github.com/guntoroyk/golang-restful-api/helper"
	"github.com/guntoroyk/golang-restful-api/model/domain"
)

type CategoryRepositoryImpl struct {
	DB *sql.DB
	Validate *validator.Validate
}

func NewCategoryRepository(DB *sql.DB, validate *validator.Validate) CategoryRepository {
	return &CategoryRepositoryImpl{
		DB: DB,
		Validate: validate,
	}
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, category domain.Category) domain.Category {
	err := repository.Validate.Struct(category)
	helper.PanicIfError(err)

	SQL := "insert into category(name) values ($1) returning id"

	tx, err := repository.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	rows, err := tx.QueryContext(ctx, SQL, category.Name)
	helper.PanicIfError(err)
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&category.Id)
		helper.PanicIfError(err)
		return category
	}

	return category
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, category domain.Category) domain.Category {
	err := repository.Validate.Struct(category)
	helper.PanicIfError(err)

	SQL := "update category set name = $1 where id = $2"

	tx, err := repository.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	_, err = tx.ExecContext(ctx, SQL, category.Name, category.Id)
	helper.PanicIfError(err)

	return category
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, category domain.Category) {
	SQL := "delete from category where id = $1"

	tx, err := repository.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	_, err = tx.ExecContext(ctx, SQL, category.Id)
	helper.PanicIfError(err)
}

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, categoryId int) (domain.Category, error) {
	SQL := "select id, name from category where id = $1"

	tx, err := repository.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	rows, err :=  tx.QueryContext(ctx, SQL, categoryId)
	helper.PanicIfError(err)
	defer rows.Close()

	category := domain.Category{}
	if rows.Next() {
		err := rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		return category, nil
	} else {
		return category, errors.New("category is not found")
	}
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context) []domain.Category {
	SQL := "select id, name from category"

	tx, err := repository.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		category := domain.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		categories = append(categories, category)
	}
	return categories
}
