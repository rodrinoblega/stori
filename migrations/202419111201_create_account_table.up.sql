CREATE TABLE IF NOT EXISTS accounts (
       account_id serial PRIMARY KEY,
       name varchar,
       mail varchar
);

INSERT INTO accounts(name, mail) VALUES ('ACCOUNT BASE', 'rnoblega@gmail.com');