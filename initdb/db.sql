CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

INSERT INTO campaigns (name) VALUES ('Первая запись');

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER NOT NULL REFERENCES campaigns (id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    priority INTEGER NOT NULL,
    removed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);