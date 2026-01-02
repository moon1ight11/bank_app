-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS bank_app;

CREATE TYPE bank_app.ROLES AS ENUM ('Basic', 'Verificator', 'Admin') ;
CREATE TYPE bank_app.CURRENCIES AS ENUM ('USD', 'EUR', 'RUB');

CREATE TABLE
    bank_app.users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR NOT NULL,
        surname VARCHAR NOT NULL,
        email VARCHAR UNIQUE NOT NULL,
        phone_number VARCHAR,
        password VARCHAR,
        timezone VARCHAR DEFAULT 'UTC',
        role bank_app.Roles NOT NULL DEFAULT 'Basic',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ
    );
CREATE INDEX idx_users_email ON bank_app.users (email);
CREATE INDEX idx_users_phone ON bank_app.users (phone_number);

CREATE TABLE
    bank_app.accounts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL REFERENCES bank_app.users (id),
        balance INTEGER DEFAULT 0,
        currency bank_app.CURRENCIES NOT NULL,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ
    );
CREATE INDEX idx_accounts_user_id ON bank_app.accounts (user_id);

CREATE TABLE
    bank_app.transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_from UUID NOT NULL REFERENCES bank_app.users (id),
        account_from UUID NOT NULL REFERENCES bank_app.accounts (id),
        user_to UUID NOT NULL REFERENCES bank_app.users (id),
        account_to UUID NOT NULL REFERENCES bank_app.accounts (id),
        amount INTEGER NOT NULL,
        currency bank_app.CURRENCIES NOT NULL,
        timestamp TIMESTAMPTZ DEFAULT now ()
    );
CREATE INDEX idx_transactions_user_from ON bank_app.transactions (user_from);
CREATE INDEX idx_transactions_account_from ON bank_app.transactions (account_from);
CREATE INDEX idx_transactions_user_to ON bank_app.transactions (user_to);
CREATE INDEX idx_transactions_account_to ON bank_app.transactions (account_to);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bank_app.transactions;
DROP TABLE IF EXISTS bank_app.accounts;
DROP TABLE IF EXISTS bank_app.users;

DROP TYPE IF EXISTS bank_app.CURRENCIES;
DROP TYPE IF EXISTS bank_app.ROLES;

DROP SCHEMA IF EXISTS bank_app;
-- +goose StatementEnd