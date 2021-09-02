package categoryview

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
)

type Consumer struct{}

type CategoryResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (cvc *Consumer) HandleMessage(message *nsq.Message) error {
	var msg struct {
		Event          string            `json:"event"`
		CategoryDetail *CategoryResponse `json:"category_detail"`
	}

	if err := json.Unmarshal(message.Body, &msg); err != nil {
		return err
	}

	log.Println("Consumer got event " + msg.Event)
	log.Println("Details: ", msg.CategoryDetail)

	return nil
}
