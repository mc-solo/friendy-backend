create table refresh_tokens (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash text unique not null,
    expires_at timestamptz not null,
    created_at timestamptz default now()

);

-- index
create index idx_refresh_tokens_user_id on refresh_tokens(user_id);