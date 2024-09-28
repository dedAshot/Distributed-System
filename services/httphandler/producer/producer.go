package producer

import (
	"encoding/json"
	"fmt"
	"httphandler/store"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var Producer *kafka.Producer
var topic string = "messages"

func New(KafkaAddr, KafkaProducerSettings string) error {

	err := connectToKafka(KafkaAddr, KafkaProducerSettings)
	for retries := 10; retries > 0 && err != nil; retries-- {
		time.Sleep(time.Second)
		err = connectToKafka(KafkaAddr, KafkaProducerSettings)
	}
	if err != nil {
		return err
	}

	return nil
}

func connectToKafka(KafkaAddr, KafkaProducerSettings string) error {

	var err error
	Producer, err = kafka.NewProducer(
		&kafka.ConfigMap{"bootstrap.servers": KafkaAddr, "request.required.acks": 1})
	if err != nil {
		return fmt.Errorf("connect to kafka error: %w", err)
	}

	go DeliveryReportHandler()

	return nil
}

// Delivery report handler for produced messages
func DeliveryReportHandler() {
	for e := range Producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		case kafka.Error:
			// Generic client instance-level errors, such as
			// broker connection failures, authentication issues, etc.
			//
			// These errors should generally be considered informational
			// as the underlying client will automatically try to
			// recover from any errors encountered, the application
			// does not need to take action on them.
			fmt.Printf("Error: %v\n", ev)
		default:
			fmt.Printf("Ignored event: %s\n", ev)

		}
	}
}

func SendMsg(msg *store.Message) error {

	data, err := json.Marshal(msg)

	if err != nil {
		return fmt.Errorf("json Marshal error: %w", err)
	}

	err = Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
	}, nil)

	// Wait for message deliveries before shutting down
	Producer.Flush(2 * 1000)

	return err
}
