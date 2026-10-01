-- +goose Up
CREATE TABLE public.users (
id uuid PRIMARY KEY,
name varchar NOT NULL,
password text NOT NULL,
created_at timestamptz NOT NULL DEFAULT now(),
CONSTRAINT users_name_unique UNIQUE (name)
);

CREATE TABLE posts (
id uuid NOT NULL,
author_id uuid NOT NULL,
text text NOT NULL,
created_at timestamptz NOT NULL DEFAULT now(),

PRIMARY KEY (id, created_at),

FOREIGN KEY (author_id) REFERENCES users(id)
) PARTITION BY RANGE (created_at);

CREATE TABLE posts_2025 PARTITION OF posts
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE posts_2026 PARTITION OF posts
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE posts_2027 PARTITION OF posts
    FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');

CREATE TABLE post_likes (
post_id uuid NOT NULL,
post_created_at timestamptz NOT NULL,
user_id uuid NOT NULL,
created_at timestamptz NOT NULL DEFAULT now(),

PRIMARY KEY (post_id, post_created_at, user_id),

FOREIGN KEY (post_id, post_created_at)
REFERENCES posts(id, created_at),

FOREIGN KEY (user_id)
REFERENCES users(id)
);

-- +goose Down
DROP TABLE IF EXISTS post_likes;
DROP TABLE IF EXISTS posts_2025;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;