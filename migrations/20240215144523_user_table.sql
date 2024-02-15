-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     surname VARCHAR(100) NOT NULL,
                                     name VARCHAR(100) NOT NULL,
                                     middlename VARCHAR(100),
                                     phone_number DECIMAL(10,0) NOT NULL,
                                     email VARCHAR(80) NOT NULL,
                                     password VARCHAR(80),
                                     sex CHAR(1),
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd