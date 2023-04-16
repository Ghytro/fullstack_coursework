-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    bio VARCHAR,
    avatar_image_id CHAR(24) UNIQUE
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL,
    caption TEXT,
    image_id CHAR(24) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX posts_user_id_idx ON posts (user_id);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX comments_post_id_user_id_idx ON comments (post_id, user_id);

CREATE TABLE likes (
    post_id UUID NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX likes_post_id_user_id_idx ON likes (post_id, user_id);

-- +goose StatemendEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS posts_user_id_idx;
DROP INDEX IF EXISTS comments_post_id_user_id_idx;
DROP INDEX IF EXISTS likes_post_id_user_id_idx;

DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

-- +goose StatemendEnd
