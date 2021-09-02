package main

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
	"time"
)

type Message struct {
	Name      string
	Content   string
	Timestamp string
}

func main() {
	config := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	//	Init topic
	topic := "Topic_Example"
	msg := Message{
		Name:      "Message name example",
		Content:   "Message content example",
		Timestamp: time.Now().String(),
	}

	//	convert message to []byte
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	//	Publish message
	err = producer.Publish(topic, payload)
	if err != nil {
		log.Println(err)
	}
}
