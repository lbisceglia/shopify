CREATE TABLE IF NOT EXISTS items (
    id CHAR(20) PRIMARY KEY,
    sku VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR,
    price_cad FLOAT,
    quantity INTEGER NOT NULL,
    date_added TIMESTAMPTZ NOT NULL,
    last_updated TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_items (
    id CHAR(20) PRIMARY KEY,
    sku VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR,
    price_cad FLOAT,
    quantity INTEGER NOT NULL,
    date_added TIMESTAMPTZ NOT NULL,
    last_updated TIMESTAMPTZ NOT NULL,
    deletion_comments TEXT,
    deleted_on TIMESTAMPTZ NOT NULL
);