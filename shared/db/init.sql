CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_url TEXT UNIQUE NOT NULL,
    long_url  TEXT NOT NULL,
    batch_id  TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);