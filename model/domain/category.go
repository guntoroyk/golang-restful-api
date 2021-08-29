package domain

type Category struct {
	Id   int
	Name string `validate:"required,max=200,min=1"`
}
