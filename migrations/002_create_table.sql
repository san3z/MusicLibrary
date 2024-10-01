-- +goose Up
CREATE TABLE public."musicLibrary"
(
    id integer NOT NULL,
    band text,
    song text,
    release_date date,
    album text,
    lyrics text,
    link text,
    PRIMARY KEY (id)
);