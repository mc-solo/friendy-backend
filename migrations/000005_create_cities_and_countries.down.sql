-- remove city_id from users first (depends on cities)
drop index if exists idx_users_city_id;
alter table users drop column if exists city_id;

-- drop tables in reverse order
drop table if exists cities;
drop table if exists countries;
