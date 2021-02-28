CREATE DATABASE IF NOT EXISTS todo_tutorial;

CREATE TABLE IF NOT EXISTS todo_tutorial.todos (
	id VARCHAR(36) NOT NULL,
	user_id VARCHAR(36) NOT NULL,
	content TEXT NOT NULL,
	completed BOOLEAN NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL,
	PRIMARY KEY (id)
);

CREATE INDEX idx_user_id ON todo_tutorial.todos(user_id);
CREATE INDEX idx_user_completed ON todo_tutorial.todos(completed);
CREATE INDEX idx_user_created_at ON todo_tutorial.todos(created_at);
CREATE INDEX idx_user_deleted ON todo_tutorial.todos(deleted);
