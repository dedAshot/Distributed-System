CREATE TABLE messages (
    id SERIAL,
    url TEXT NOT NULL,
	req_body BYTEA,
    PRIMARY KEY(id)
);