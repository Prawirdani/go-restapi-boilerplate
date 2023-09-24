CREATE TABLE IF NOT EXISTS users (
    id varchar(255) NOT NULL PRIMARY_KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp
)