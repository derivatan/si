CREATE
    EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE artists
(
    id         UUID      NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    name       TEXT
);

CREATE TABLE contacts
(
    id         UUID      NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    email      TEXT,
    phone      TEXT,
    artist_id  UUID      NOT NULL
);

CREATE TABLE albums
(
    id         UUID      NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    name       TEXT,
    year       INT,
    artist_id  UUID      NOT NULL
);
