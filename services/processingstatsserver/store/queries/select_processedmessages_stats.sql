SELECT id, datep, topic, partition, partitionOffset
FROM processed_messages
    WHERE processed_messages.id > $1
    ORDER BY id DESC
    LIMIT $2;