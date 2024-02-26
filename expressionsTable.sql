CREATE TABLE expressions (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    status TEXT,
    result TEXT,
    created_at TIMESTAMP NOT NULL,
    calculated_at TIMESTAMP
);
