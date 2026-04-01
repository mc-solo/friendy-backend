create table user_preferences(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null unique references users(id) on delete cascade,

    created_at timestamptz default now(),
    updated_at timestamptz default now()
);
