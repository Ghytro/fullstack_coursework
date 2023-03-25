-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    first_name VARCHAR,
    last_name VARCHAR,
    password VARCHAR NOT NULL,
    bio VARCHAR,
    avatar_image_id UUID UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    deactivated BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX users_username_idx ON users (username);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    caption TEXT,
    image_id UUID UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    likes INT NOT NULL DEFAULT 0,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX posts_user_id_idx ON posts (user_id);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX comments_post_id_idx ON comments (post_id);
CREATE INDEX comments_user_id_idx ON comments (user_id);

CREATE TABLE subscriptions (
    subscriber_id INT NOT NULL,
    publisher_id INT NOT NULL,

    FOREIGN KEY (subscriber_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (publisher_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX subscriptions_subscriber_id ON subscriptions (subscriber_id);
CREATE INDEX subscriptions_publisher_id ON subscriptions (publisher_id);

-- +goose StatemendEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS users_username_idx;
DROP INDEX IF EXISTS posts_user_id_idx;
DROP INDEX IF EXISTS comments_post_id_idx;
DROP INDEX IF EXISTS comments_user_id_idx;
DROP INDEX IF EXISTS subscriptions_subscriber_id;
DROP INDEX IF EXISTS subscriptions_publisher_id;

-- +goose StatemendEnd
