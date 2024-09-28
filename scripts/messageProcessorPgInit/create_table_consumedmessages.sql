CREATE TABLE processed_messages (
    id SERIAL,
    datep TIMESTAMP,
    topic TEXT,
    partition INT,
    partitionOffset INT,
    PRIMARY KEY(id)
);