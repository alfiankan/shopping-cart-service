CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS cart (
    id uuid DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);
