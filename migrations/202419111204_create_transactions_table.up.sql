CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(255) NOT NULL,
    transactions_type VARCHAR(255) NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    amount NUMERIC NOT NULL,
    account_id integer NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts (account_id)
);