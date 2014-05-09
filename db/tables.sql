-- Copyright (C) 2014 Miquel Sabaté Solà
-- This file is licensed under the MIT license.
-- See the LICENSE file.


create table users (
    id uuid primary key,
    name varchar(255) unique not null,
    auth_token varchar(255) unique,
    password_hash text,
    created_at timestamp
);

create table players (
    id uuid primary key,
    name varchar(255) unique not null,
    created_at timestamp
);
