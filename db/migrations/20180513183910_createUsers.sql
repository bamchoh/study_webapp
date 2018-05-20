
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
	id TEXT NOT NULL,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	PRIMARY KEY (id),
	UNIQUE(email)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS pgcrypto;

