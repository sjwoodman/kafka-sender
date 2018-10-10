package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/matzew/kafka-receiver/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		log.Print("Received a POST request.")

		cfg := sarama.NewConfig()
		cfg.Producer.RequiredAcks = sarama.WaitForAll
		cfg.Producer.Retry.Max = 5
		cfg.Producer.Return.Successes = true

		flowConfig := config.GetConfig()

		producer, err := sarama.NewSyncProducer([]string{flowConfig.BootStrapServers}, cfg)
		if err != nil {
			// Should not reach here
			panic(err)
		}

		defer func() {
			if err := producer.Close(); err != nil {
				// Should not reach here
				panic(err)
			}
		}()

		bodyBuffer, _ := ioutil.ReadAll(r.Body)

		msg := &sarama.ProducerMessage{
			Topic: flowConfig.KafkaTopic,
			Value: sarama.StringEncoder(bodyBuffer),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", flowConfig.KafkaTopic, partition, offset)
	}
}

func main() {
	flag.Parse()
	log.Print("CloudEvent Kafka Sender started.")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
