CREATE TABLE processed_messages (
    id SERIAL,
    datep DATE,
    topic TEXT,
    partition INT,
    partitionOffset INT,
    PRIMARY KEY(id)
);

CREATE INDEX processed_messages_idx ON processed_messages (id);