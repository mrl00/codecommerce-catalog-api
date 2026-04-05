-- Initialization schema for codecommerceapi

CREATE TABLE IF NOT EXISTS tb_category (
    pk_category UUID PRIMARY KEY,
    tx_name TEXT NOT NULL,
    ts_category_created_at TIMESTAMP NOT NULL,
    ts_category_updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tb_product (
    pk_product UUID PRIMARY KEY,
    tx_name TEXT NOT NULL,
    tx_description TEXT NOT NULL DEFAULT '',
    nr_price DOUBLE PRECISION NOT NULL,
    tx_image_url TEXT NOT NULL DEFAULT '',
    fk_category UUID NOT NULL REFERENCES tb_category(pk_category) ON DELETE CASCADE,
    ts_product_created_at TIMESTAMP NOT NULL,
    ts_product_updated_at TIMESTAMP NOT NULL
);
