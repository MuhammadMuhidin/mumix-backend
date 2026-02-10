create table if not exists todos (
  id uuid primary key default gen_random_uuid(),
  title text not null,
  completed boolean default false,
  created_at timestamptz default now()
);
