-- countries table
create table countries (
    id uuid primary key default gen_random_uuid(),
    name varchar(255) unique not null,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- cities table (depends on countries)
create table cities (
    id uuid primary key default gen_random_uuid(),
    name varchar(255) not null,
    country_id uuid not null references countries(id) on delete cascade,
    latitude double precision,
    longitude double precision,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint idx_city_country_name unique (name, country_id)
);

-- add city_id column to users
alter table users add column city_id uuid references cities(id) on delete set null;
create index idx_users_city_id on users(city_id);
