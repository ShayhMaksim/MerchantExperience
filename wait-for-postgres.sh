#!/bin/sh
# wait-for-postgres.sh

set -e
  
shift
cmd="$@"

# echo "Настраивается вручную БД: смотреть файлы processing/database.go"
# PGPASSWORD=$DB_PASSWORD psql -h "database" -U "postgres" <<-EOSQL
#   DROP DATABASE base;
#   DROP DATABASE base;
#   DROP View responsibility;
#   DROP Table sellers;
#   DROP Table products;
# EOSQL

PGPASSWORD=$DB_PASSWORD psql -h "database" -U "postgres" <<-EOSQL
CREATE DATABASE base;

\c base

CREATE TABLE products
(
    offer_id SERIAL PRIMARY KEY,
    name CHARACTER VARYING(30),
    price NUMERIC(16,2),
    quantity INTEGER,
    available BOOLEAN
);

CREATE TABLE sellers
(
	seller_id SERIAL, 
    offer_id SERIAL
);

CREATE VIEW responsibility AS
  SELECT sellers.seller_id, products.offer_id,products.name,products.price,products.quantity
  FROM sellers,products
  WHERE products.offer_id=sellers.offer_id;
EOSQL

#docker exec -it 8ed5f20459e4 psql -U "postgres"


echo "Postgres is up - executing command"
exec $cmd