package messaging

import (
	"github.com/nsqio/go-nsq"
	"log"
)

type (
	ConsumerConfig struct {
		Topic         string
		Channel       string
		LookupAddress string
		MaxAttempts   uint16
		MaxInFlight   int
		Handler       nsq.HandlerFunc
	}

	Consumer struct {
		consumer      *nsq.Consumer
		lookupAddress string
		handler       nsq.HandlerFunc
	}
)

func NewConsumer(cfg ConsumerConfig) Consumer {
	nsqConf := nsq.NewConfig()
	nsqConf.MaxAttempts = cfg.MaxAttempts
	nsqConf.MaxInFlight = cfg.MaxInFlight

	topic := cfg.Topic
	c, err := nsq.NewConsumer(topic, cfg.Channel, nsqConf)
	if err != nil {
		log.Fatal(err)
	}

	return Consumer{
		consumer:      c,
		lookupAddress: cfg.LookupAddress,
		handler:       cfg.Handler,
	}
}

func (c *Consumer) Run() {
	c.consumer.AddHandler(c.handler)
	err := c.consumer.ConnectToNSQLookupd(c.lookupAddress)
	if err != nil {
		log.Fatal(err)
	}
}
