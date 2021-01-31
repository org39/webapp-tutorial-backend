CREATE DATABASE IF NOT EXISTS test;

CREATE TABLE test.users (
	id VARCHAR(36) NOT NULL,
	email VARCHAR(256) NOT NULL,
	password VARCHAR(64) NOT NULL,
	PRIMARY KEY (id)
);
