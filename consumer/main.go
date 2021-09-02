package main

import (
	"github.com/guntoroyk/golang-restful-api/consumer/messaging"
	"github.com/guntoroyk/golang-restful-api/consumer/messaging/consumer/categoryview"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cvc := categoryview.NewConsumer()

	consumerCfg := messaging.ConsumerConfig{
		Topic:         "category_view",
		Channel:       "golang_restful_api_consumer",
		LookupAddress: "127.0.0.1:4161",
		MaxAttempts:   10,
		MaxInFlight:   100,
		Handler:       cvc.HandleMessage,
	}

	consumerCategoryView := messaging.NewConsumer(consumerCfg)

	consumerCategoryView.Run()

	// keep app alive until terminated
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("Application terminated")
	}
}
