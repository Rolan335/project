-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS blogs(
    id uuid PRIMARY KEY NOT NULL,
    users_id UUID NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (users_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS posts(
    id uuid PRIMARY KEY NOT NULL,
    blogs_id UUID NOT NULL,
    title TEXT NOT NULL,
    "text" TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (blogs_id) REFERENCES blogs(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS blogs;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
