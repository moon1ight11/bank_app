-- +goose Up
-- +goose StatementBegin
INSERT INTO 
    bank_app.users (id, name, surname, email, phone_number, password, role)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'admin', 'zero', 'nulladmin@mail.com', '+70000000000', '$2a$10$rQxxz3Hg/g2IRyIhm8K9/u.q0onCMAve1bGcfbHNeBYbvojhjjKVS', 'Admin');

INSERT INTO
    bank_app.accounts (id, user_id, currency)
VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'RUB'), 
    ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'EUR'), 
    ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'USD');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM bank_app.accounts 
WHERE id IN (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000002', 
    '00000000-0000-0000-0000-000000000003'
);
DELETE FROM bank_app.users 
WHERE id = '00000000-0000-0000-0000-000000000001';
-- +goose StatementEnd