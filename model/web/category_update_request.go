package web

type CategoryUpdateRequest struct {
	Id   int
	Name string `json:"name"`
}
