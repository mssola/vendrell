-- Copyright (C) 2014-2015 Miquel Sabaté Solà
-- This file is licensed under the MIT license.
-- See the LICENSE file.


create table users (
    id uuid primary key,
    name varchar(255) unique not null check (name <> ''),
    password_hash text,
    created_at timestamp
);

create table players (
    id uuid primary key,
    name varchar(255) unique not null check (name <> ''),
    created_at timestamp
);

create table ratings (
    id serial primary key,
    value int not null,
    player_id uuid references players(id) on delete cascade,
    created_at timestamp
);

