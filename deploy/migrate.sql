DROP TABLE IF EXISTS users

CREATE TABLE users (
	id serial PRIMARY KEY, 
	name VARCHAR (50) UNIQUE NOT NULL,
	password_hash VARCHAR NOT NULL
)

