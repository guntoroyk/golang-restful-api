package main

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Message struct {
	Name      string
	Content   string
	Timestamp string
}

type MessageHandler struct{}

func (h MessageHandler) HandleMessage(m *nsq.Message) error {
	var request Message

	if err := json.Unmarshal(m.Body, &request); err != nil {
		log.Println("error when unmarshalling the message body, err: ", err)
		return err
	}

	//Print the Message
	log.Println("Message")
	log.Println("--------------------")
	log.Println("Name : ", request.Name)
	log.Println("Content : ", request.Content)
	log.Println("Timestamp : ", request.Timestamp)
	log.Println("--------------------")
	log.Println("")
	// Will automatically set the message as finish
	return nil
}

func main() {
	config := nsq.NewConfig()

	config.MaxAttempts = 10

	config.MaxInFlight = 5

	config.MaxRequeueDelay = time.Second * 900
	config.DefaultRequeueDelay = time.Second * 0

	topic := "Topic_Example"
	channel := "Channel_Example"

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(&MessageHandler{})

	err = consumer.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		log.Fatalf("failed to connect to NSQLookupd, err: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
}
