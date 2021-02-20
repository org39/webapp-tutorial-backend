CREATE DATABASE IF NOT EXISTS todo_tutorial;

CREATE TABLE IF NOT EXISTS todo_tutorial.users (
	id VARCHAR(36) NOT NULL,
	email VARCHAR(256) NOT NULL,
	password VARCHAR(64) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_user_email ON todo_tutorial.users(email);
