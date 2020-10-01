CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
	email text NOT NULL UNIQUE,
	password text NOT NULL,
	name text NOT NULL
);

INSERT INTO users (email, password, name) VALUES
('admin@gmail.com', 'cGFzc3dvcmQ=', 'admin');