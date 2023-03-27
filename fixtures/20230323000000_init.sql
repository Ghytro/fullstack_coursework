-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    first_name VARCHAR,
    last_name VARCHAR,
    password VARCHAR NOT NULL,
    bio VARCHAR,
    avatar_image_id CHAR(24) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    deactivated BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL,
    caption TEXT,
    image_id CHAR(24) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

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

CREATE UNIQUE INDEX comments_post_id_user_id_idx ON comments (post_id, user_id);

CREATE TABLE likes (
    post_id UUID NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX likes_post_id_user_id_idx ON likes (post_id, user_id);

CREATE TABLE comment_likes (
    comment_id UUID NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX comment_likes_comment_id_user_id_idx ON comment_likes (comment_id, user_id);

CREATE TABLE subscriptions (
    subscriber_id INT NOT NULL,
    publisher_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (subscriber_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (publisher_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX subscriptions_subscriber_id_publisher_id_idx ON subscriptions (subscriber_id, publisher_id);

-- +goose StatemendEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS likes_post_id_user_id_idx;
DROP INDEX IF EXISTS subscriptions_subscriber_id_publisher_id_idx;
DROP INDEX IF EXISTS comment_likes_comment_id_user_id_idx;
DROP INDEX IF EXISTS posts_user_id_idx;
DROP INDEX IF EXISTS comments_post_id_user_id_idx;

DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS comment_likes;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;

-- +goose StatemendEnd
