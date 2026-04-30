-- +bun:up
CREATE TABLE IF NOT EXISTS stocks (
    id       TEXT NOT NULL,
    name     TEXT NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS products (
    id       TEXT    NOT NULL,
    name     TEXT    NOT NULL,
    qtd      INTEGER NOT NULL DEFAULT 0,
    price    INTEGER NOT NULL,
    stock_id TEXT    NOT NULL,

    PRIMARY KEY (id),

    CONSTRAINT fk_products_stock
        FOREIGN KEY (stock_id)
        REFERENCES stocks (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
