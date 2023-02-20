-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL SECOND(0) NOT NULL,
    description TEXT,
    userid UUID NOT NULL,
    notify INTERVAL SECOND(0),
    sent BOOLEAN NOT NULL DEFAULT FALSE,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd