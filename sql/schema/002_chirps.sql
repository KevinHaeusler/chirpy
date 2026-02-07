-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    body text NOT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY (user_id) references users(id)
    ON DELETE CASCADE);


-- +goose Down
DROP TABLE chirps;