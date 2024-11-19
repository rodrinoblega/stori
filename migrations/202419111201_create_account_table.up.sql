CREATE TABLE IF NOT EXISTS accounts (
       account_id serial PRIMARY KEY,
       name varchar
);

INSERT INTO accounts(name) VALUES ('ACCOUNT BASE');