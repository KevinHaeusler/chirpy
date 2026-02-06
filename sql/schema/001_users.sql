-- +goose Up
CREATE TABLE users (
   id uuid PRIMARY KEY,
   created_at timestamp NOT NULL DEFAULT now(),
   updated_at timestamp NOT NULL DEFAULT now(),
   email TEXT NOT NULL UNIQUE

);

-- +goose Down
DROP TABLE users;