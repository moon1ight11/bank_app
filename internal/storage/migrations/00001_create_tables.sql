-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS bank_app;

CREATE TYPE bank_app.OPERATION_TYPES AS ENUM ('incoming', 'outlay', 'transfer');
CREATE TYPE bank_app.CURRENCIES AS ENUM ('USD', 'EUR', 'RUB');

CREATE TABLE
    bank_app.users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR NOT NULL,
        surname VARCHAR NOT NULL,
        email VARCHAR UNIQUE NOT NULL,
        phone_number VARCHAR,
        timezone VARCHAR DEFAULT 'UTC',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ
    );
CREATE INDEX idx_users_email ON bank_app.users(email);
CREATE INDEX idx_users_phone ON bank_app.users(phone_number);

CREATE TABLE
    bank_app.accounts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL REFERENCES bank_app.users (id),
        balance INTEGER DEFAULT 0,
        currency bank_app.CURRENCIES NOT NULL DEFAULT 'RUB',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ
    );
CREATE INDEX idx_accounts_user_id ON bank_app.accounts(user_id);


CREATE TABLE
    bank_app.operations (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL REFERENCES bank_app.users (id),
        account_id UUID NOT NULL REFERENCES bank_app.accounts(id),
        type bank_app.OPERATION_TYPES NOT NULL,
        amount INTEGER NOT NULL,
        timestamp TIMESTAMPTZ DEFAULT now()
    );
CREATE INDEX idx_operations_user_id ON bank_app.operations(user_id);
CREATE INDEX idx_operations_account_id ON bank_app.operations(account_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bank_app.operations;
DROP TABLE IF EXISTS bank_app.accounts;
DROP TABLE IF EXISTS bank_app.users;

DROP TYPE IF EXISTS bank_app.CURRENCIES;
DROP TYPE IF EXISTS bank_app.OPERATION_TYPES;

DROP SCHEMA IF EXISTS bank_app;
-- +goose StatementEnd