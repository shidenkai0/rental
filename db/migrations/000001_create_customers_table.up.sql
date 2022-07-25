-- Create customers table with id, name columns
BEGIN;
CREATE TABLE customers (
    "id" SERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL
);
COMMIT;
