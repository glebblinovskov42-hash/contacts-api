CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    number VARCHAR(20) NOT NULL,
    is_favourite BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() 
);