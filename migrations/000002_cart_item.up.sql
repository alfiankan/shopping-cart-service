CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS cart_item (
    id uuid DEFAULT uuid_generate_v4(),
    product_code uuid,
    product_name TEXT NOT NULL,
    qty INT NOT NULL DEFAULT 1,
    cart_id uuid NOT NULL,
    CONSTRAINT fk_cart_id FOREIGN KEY (cart_id) REFERENCES cart (id)
);
