-- +goose Up
-- +goose StatementBegin
INSERT INTO 
    bank_app.users (id, name, surname, email, phone_number, password, role)
VALUES
    ('00000000-0000-0000-0000-000000000000', 'admin', 'zero', 'nulladmin@mail.com', '+70000000000', '112900', 'Admin');

INSERT INTO
    bank_app.accounts (id, user_id, currency)
VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000000', 'RUB'), 
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000000', 'EUR'), 
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000000', 'USD');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM bank_app.accounts;
DELETE FROM bank_app.users;
-- +goose StatementEnd