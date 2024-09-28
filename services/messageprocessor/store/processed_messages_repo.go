package store

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var ProcessedMessageRepository messageRepository

type StatRow struct {
	Timestamp       time.Time ``
	TopicName       string    ``
	Partition       int32     ``
	PartitionOffset string    ``
}

type messageRepository struct{}

func (m *messageRepository) SaveProcessedRow(msg *kafka.Message) error {
	sendConsMsgQueryTempl := queryTemplates["insert_message.sql"]
	fmt.Println(msg.Timestamp.Format(time.RFC3339))
	//time, _ := time.Parse(time.RFC3339, msg.Timestamp)
	_, err := db.Exec(
		sendConsMsgQueryTempl,
		msg.Timestamp.Format(time.RFC3339),
		msg.TopicPartition.Topic,
		msg.TopicPartition.Partition,
		msg.TopicPartition.Offset.String(),
	)

	if err != nil {
		return fmt.Errorf("message saving in db error: %w", err)
	}

	return nil
}
