SELECT id, datep, topic, partition, partitionOffset
FROM processed_messages
    ORDER BY id DESC
    LIMIT $1;