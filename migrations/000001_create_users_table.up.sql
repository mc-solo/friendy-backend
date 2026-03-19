create extension if not exists pgcrypto;

create table users (
    id uuid primary key default gen_random_uuid(),
    email varchar(255) unique not null,
    username varchar(25) unique,
    password_hash text not null,
    first_name varchar(100),
    last_name varchar(100),

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_users_deleted_at on users(deleted_at);