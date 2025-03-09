-- +migrate Up
CREATE TABLE messages (
                       id_messages SERIAL PRIMARY KEY,
                       messages TEXT NOT NULL,
                       time_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);
-- +migrate Down
DROP TABLE messages;