-- CREATE cars table with id, customer_id foreign key, make, model, year columns
BEGIN;
CREATE TABLE cars (
    id serial PRIMARY KEY,
    customer_id integer,
    make varchar(255) NOT NULL,
    model varchar(255) NOT NULL,
    year integer NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers (id)
);
COMMIT;
