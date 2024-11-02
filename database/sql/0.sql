create table migrations (
    id integer constraint pk_migrations primary key,
    applied_at date not null default now()
);

create table platforms (
    id uuid primary key default gen_random_uuid(),
    name varchar(200) not null constraint ix_platform_name unique,
    short_name varchar(5) not null constraint ix_platform_short_name unique
);

create table games (
    id uuid primary key default gen_random_uuid(),
    title varchar(500) not null,
    platform_id uuid not null references platforms(id) on delete restrict,
    owned boolean not null default false,
    release_date timestamp with time zone null,
    released boolean not null
);

create table playthroughs (
    id uuid primary key default gen_random_uuid(),
    game_id uuid not null references games(id) on delete cascade,
    start_date timestamp with time zone not null,
    end_date timestamp with time zone null,
    -- status values:
    -- 0 - in progress
    -- 1 - completed
    -- 2 - dropped
    -- 3 - retired (for endless games)
    -- 4 - suspended (should not be set with end_date)
    status smallint not null default 0 check ( status >= 0 and status <= 4 ),
    runtime_minutes integer null default null
);
