package messaging

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"log"
)

type (
	ProducerConfig struct {
		NsqdAddress string
	}

	Producer struct {
		prod *nsq.Producer
	}
)

func NewProducer(cfg ProducerConfig) Producer {
	p, err := nsq.NewProducer(cfg.NsqdAddress, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	return Producer{prod: p}
}

func (p *Producer) Publish(topic string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.prod.Publish(topic, payload)
}
