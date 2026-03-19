create extension if not exists pgcrypto;

-- enum defn //the below are also defined on our enums.go file
create type gender_type as enum ('male', 'female');
create type body_type as enum('slim', 'athletic', 'chubby', 'curvy', 'fit', 'other');
create type educational_level as enum ('high_school', 'bachelor', 'master', 'phd', 'university_drop_out', 'high_school_drop_out', 'home_schooled', 'other');
create type language as enum ('am','en', 'or', 'tg', 'gz', 'sp', 'fr', 'it', 'other');



create table user_profile(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null unique references users(id) on delete cascade,
    gender gender_type,
    bio varchar(255),
    body_type body_type,
    birthdate date,
    height_cm decimal(5,2),
    lang language not null default 'en',
    edu_level educational_level,
    profile_completion_rate integer not null default 0,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- indexes 
create index idx_user_profile_gender on user_profile(gender);
create index idx_user_profile_birth_date on user_profile(birthdate);
create index idx_user_profile_lang on user_profile(lang);
