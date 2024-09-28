package store

import (
	"time"
)

var ProcessedMessageRepository messageRepository

type StatRow struct {
	Id              int       `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	TopicName       string    `json:"topicname"`
	Partition       int32     `json:"partition"`
	PartitionOffset string    `json:"partitionoffset"`
}
type messageRepository struct{}

func GetMessages(startId, count int) ([]*StatRow, error) {

	query := queryTemplates["select_processedmessages_stats.sql"]

	resRows := make([]*StatRow, 0, count)

	rows, err := db.Query(query, startId, count)
	if err != nil {
		return nil, err
	}

	for i := 0; rows.Next(); i++ {
		row := &StatRow{}
		err := rows.Scan(&row.Id, &row.Timestamp, &row.TopicName, &row.Partition, &row.PartitionOffset)
		if err != nil {
			return nil, err
		}
		resRows = append(resRows, row)
	}

	return resRows, nil
}

func GetLastMessages(count int) ([]*StatRow, error) {

	query := queryTemplates["select_last_processedmessages_stats.sql"]

	resRows := make([]*StatRow, 0, count)

	rows, err := db.Query(query, count)
	if err != nil {
		return nil, err
	}

	for i := 0; rows.Next(); i++ {
		row := &StatRow{}
		err := rows.Scan(&row.Id, &row.Timestamp, &row.TopicName, &row.Partition, &row.PartitionOffset)
		if err != nil {
			return nil, err
		}
		resRows = append(resRows, row)
	}

	return resRows, nil
}
