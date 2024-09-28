package consumer

import (
	"fmt"
	"messageprocessor/store"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var Consumer *kafka.Consumer
var topic string = "messages"
var runSignal *bool = new(bool)

func New(KafkaAddr, KafkaConsumerGroup string) error {

	err := connectToKafka(KafkaAddr, KafkaConsumerGroup)
	for retries := 10; retries > 0 && err != nil; retries-- {
		time.Sleep(time.Second)
		err = connectToKafka(KafkaAddr, KafkaConsumerGroup)
	}
	if err != nil {
		return err
	}

	return nil
}

func connectToKafka(KafkaAddr, KafkaConsumerGroup string) error {

	var err error
	Consumer, err = kafka.NewConsumer(
		&kafka.ConfigMap{"bootstrap.servers": KafkaAddr, "group.id": KafkaConsumerGroup, "auto.offset.reset": "latest"})

	if err != nil {
		return fmt.Errorf("connect to kafka error: %w", err)
	}

	err = Consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	*runSignal = true
	consumeMessages()

	return nil
}

func consumeMessages() {
	fmt.Println("start consuming")
	for *runSignal {

		if msg, err := Consumer.ReadMessage(time.Second); err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			go func() {
				err := store.ProcessedMessageRepository.SaveProcessedRow(msg)
				if err != nil {
					fmt.Println(err)
				}
			}()

		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
