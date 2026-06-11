CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS customers (
    id    UUID         PRIMARY KEY DEFAULT uuid_generate_v4(),
    name  VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id    UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    name  VARCHAR(255)   NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0)
);

CREATE TABLE IF NOT EXISTS orders (
    id          UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID        NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    order_date  DATE        NOT NULL DEFAULT CURRENT_DATE,
    status      VARCHAR(50) NOT NULL DEFAULT 'pending'
                            CHECK (status IN ('pending','confirmed','shipped','delivered','cancelled'))
);

CREATE TABLE IF NOT EXISTS order_items (
    id         UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id   UUID           NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID           NOT NULL REFERENCES products(id),
    quantity   INTEGER        NOT NULL CHECK (quantity > 0),
    price      NUMERIC(10, 2) NOT NULL CHECK (price >= 0)
);
