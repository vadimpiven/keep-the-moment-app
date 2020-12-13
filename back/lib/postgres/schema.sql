CREATE SEQUENCE IF NOT EXISTS user_id_seq
    INCREMENT BY 1
    MINVALUE 10000000
    MAXVALUE 99999999
    START WITH 10000000
    NO CYCLE;

CREATE TABLE IF NOT EXISTS images (
    path TEXT UNIQUE,
    uploaded TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (path, uploaded)
);

INSERT INTO images (path) VALUES ('placeholder.png');

CREATE TABLE IF NOT EXISTS users (
    email TEXT UNIQUE,
    id TEXT NOT NULL UNIQUE DEFAULT 'id'||nextval('user_id_seq')::text,
    username TEXT,
    bio TEXT,
    hashtags TEXT[] NOT NULL DEFAULT '{}'::text[],
    image TEXT NOT NULL DEFAULT 'placeholder.png',
    birth DATE,
    registered TIMESTAMPTZ DEFAULT now(),
    updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (email, registered),
    FOREIGN KEY (image) REFERENCES images (path) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS hashtags (
    name TEXT UNIQUE,
    counter BIGINT NOT NULL DEFAULT '0'::bigint,
    PRIMARY KEY (name)
);

CREATE TABLE IF NOT EXISTS locations (
    email TEXT UNIQUE,
    latitude DOUBLE PRECISION NOT NULL DEFAULT '0'::double precision,
    longitude DOUBLE PRECISION NOT NULL DEFAULT '0'::double precision,
    updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (email),
    FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL UNIQUE,
    email TEXT NOT NULL,
    background INTEGER[] NOT NULL DEFAULT '{}'::integer[],
    content TEXT,
    hashtags TEXT[] NOT NULL DEFAULT '{}'::text[],
    image_1 TEXT,
    image_2 TEXT,
    image_3 TEXT,
    image_4 TEXT,
    image_5 TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT now(),
    hidden_at TIMESTAMPTZ,
    likes BIGINT NOT NULL DEFAULT '0'::bigint,
    PRIMARY KEY (id),
    FOREIGN KEY (email) REFERENCES users (email) ON DELETE RESTRICT,
    FOREIGN KEY (image_1) REFERENCES images (path) ON DELETE RESTRICT,
    FOREIGN KEY (image_2) REFERENCES images (path) ON DELETE RESTRICT,
    FOREIGN KEY (image_3) REFERENCES images (path) ON DELETE RESTRICT,
    FOREIGN KEY (image_4) REFERENCES images (path) ON DELETE RESTRICT,
    FOREIGN KEY (image_5) REFERENCES images (path) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS post_likes (
    post_id BIGINT NOT NULL,
    email TEXT NOT NULL,
    liked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE RESTRICT,
    FOREIGN KEY (email) REFERENCES users (email) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS post_comments (
    id BIGSERIAL UNIQUE,
    post_id BIGINT NOT NULL,
    email TEXT NOT NULL,
    content TEXT,
    commented_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE RESTRICT,
    FOREIGN KEY (email) REFERENCES users (email) ON DELETE RESTRICT
);
