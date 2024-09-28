SELECT * FROM processed_messages
    WHERE id > $1
    ORDER BY id ASC
    LIMIT 10;