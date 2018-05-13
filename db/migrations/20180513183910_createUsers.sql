
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION pgcrypto;

CREATE TABLE users (
	id TEXT NOT NULL,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	PRIMARY KEY (id),
	UNIQUE(email)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE users;

DROP EXTENSION pgcrypto;

