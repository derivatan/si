CREATE
    EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE artists
(
    id         UUID      NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    name       TEXT,
    year       INT
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



INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('78735385-f9aa-4121-a94c-c23589dbacc6', '2023-06-08 14:35:01.388828', NULL, NULL, 'Dummy', 1994,
        '38bfa35d-d3d8-4311-8f72-6f462b1b158b');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('446b8844-67da-482f-88c1-2bed76f541f3', '2023-06-08 14:35:01.39405', NULL, NULL, 'Portishead', 1997,
        '38bfa35d-d3d8-4311-8f72-6f462b1b158b');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('80dea8c0-c0a3-4b30-b05a-eb318ed4ff97', '2023-06-08 14:35:01.397653', NULL, NULL, 'Third', 2008,
        '38bfa35d-d3d8-4311-8f72-6f462b1b158b');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('5f45820e-b5e8-413c-833a-fbef0aff9070', '2023-06-08 14:35:01.40148', NULL, NULL, 'Meddle', 1971,
        'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('4c9a1578-98c1-49e3-805f-25f8992555d9', '2023-06-08 14:35:01.406028', NULL, NULL, 'The Dark Side of the Moon',
        1973, 'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('4f97fbfb-e7f2-415e-b716-38bc10e9ebc8', '2023-06-08 14:35:01.409946', NULL, NULL, 'Wish You Were Here', 1975,
        'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('17617313-77db-4a18-a23a-dc99a882158f', '2023-06-08 14:35:01.413591', NULL, NULL, 'The Wall', 1979,
        'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('6157c144-b4a2-4b7f-ad09-c06c5ab8742a', '2023-06-08 14:35:01.417516', NULL, NULL, 'A Momentary Lapse of Reason',
        1987, 'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('036812c9-7fc2-4d1a-943c-db570e37197e', '2023-06-08 14:35:01.420953', NULL, NULL, 'The Division Bell', 1994,
        'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('14eecfd3-bec9-42c1-b555-9ede5d2ebf91', '2023-06-08 14:35:01.424358', NULL, NULL, 'Greatest Hits', 1981,
        'a9b859d2-a905-4049-a724-8d54a0f7c4cb');
INSERT INTO albums (id, created_at, updated_at, deleted_at, name, year, artist_id)
VALUES ('e3377ea0-c57c-4a72-b3e5-e48f44b9ddd0', '2023-06-08 14:35:01.427856', NULL, NULL, 'Greatest Hits II', 1991,
        'a9b859d2-a905-4049-a724-8d54a0f7c4cb');


--
-- Data for Name: artists; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO artists (id, created_at, updated_at, deleted_at, name, year)
VALUES ('eae7150e-8ee9-4291-a2aa-1d2367aedc15', '2023-06-08 14:35:01.363225', NULL, NULL, 'Pink Floyd', 1965);
INSERT INTO artists (id, created_at, updated_at, deleted_at, name, year)
VALUES ('a9b859d2-a905-4049-a724-8d54a0f7c4cb', '2023-06-08 14:35:01.368397', NULL, NULL, 'Queen', 1970);
INSERT INTO artists (id, created_at, updated_at, deleted_at, name, year)
VALUES ('38bfa35d-d3d8-4311-8f72-6f462b1b158b', '2023-06-08 14:35:01.372094', NULL, NULL, 'Portishead', 1991);


--
-- Data for Name: contacts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO contacts (id, created_at, updated_at, deleted_at, email, phone, artist_id)
VALUES ('fd32f32d-ef51-4bd1-ac30-c74a3be0ce74', '2023-06-08 14:35:01.376174', NULL, NULL, 'pink@floyd.com', '123',
        'eae7150e-8ee9-4291-a2aa-1d2367aedc15');
INSERT INTO contacts (id, created_at, updated_at, deleted_at, email, phone, artist_id)
VALUES ('3e3b1c3e-c9af-4360-87b4-7a6dcbaab876', '2023-06-08 14:35:01.381222', NULL, NULL, 'pink@floyd.com', '456',
        'a9b859d2-a905-4049-a724-8d54a0f7c4cb');
INSERT INTO contacts (id, created_at, updated_at, deleted_at, email, phone, artist_id)
VALUES ('ddef9e45-a182-45f0-b39a-1360786542d0', '2023-06-08 14:35:01.38512', NULL, NULL, 'Queen', '789',
        '38bfa35d-d3d8-4311-8f72-6f462b1b158b');


--
-- PostgreSQL database dump complete
--

