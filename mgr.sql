-- +migrate Up
CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          message TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE messages;
