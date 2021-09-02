package web

type CategoryProducerMessage struct {
	Event          string            `json:"event"`
	CategoryDetail *CategoryResponse `json:"category_detail"`
}
