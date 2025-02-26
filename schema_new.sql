create extension if not exists pgcrypto;
-- таблица с пользователями
create table if not exists users (
  id            bigserial primary key,
  login         text not null unique,
  password_hash text not null,
  created_at    timestamptz not null default now()
);
-- -- таблица сессий
-- create table if not exists sessions (
--   id         text primary key default encode(gen_random_bytes(16),'hex'),
--   uid        bigint not null unique,  -- Уникальная сессия на пользователя
--   ip_address inet not null,           -- Храним IP-адрес
--   created_at timestamptz not null default now(),
--   expires_at timestamptz not null default now() + interval '24 hours' -- Сессия истекает через 24 часа
-- );

create table if not exists sessions (
  id         text primary key default encode(gen_random_bytes(16),'hex'),
  uid        bigint not null,       -- user id
  ip_address text not null,         -- IP адрес
  expires_at timestamptz not null,  -- Время окончания сессии
  created_at timestamptz not null default now()
);
-- таблица с файлами
create table if not exists assets (
  name       text not null,
  uid        bigint not null,
  data       bytea not null,
  created_at timestamptz not null default now(),
  primary key (name, uid)
);
-- тестовый пользователь
insert into users
(login, password_hash)
values
('alice', encode(digest('secret', 'md5'),'hex'))
on conflict do nothing;