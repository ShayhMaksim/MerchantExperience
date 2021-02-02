CREATE TABLE products
(
    offer_id SERIAL PRIMARY KEY,
    name CHARACTER VARYING(30),
    price NUMERIC(16,2),
    quantity INTEGER,
    available BOOLEAN
);