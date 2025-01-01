create table migrations (
    id integer constraint pk_migrations primary key,
    applied_at date not null default now()
);

create table users (
    id uuid primary key default gen_random_uuid(),
    provider_id varchar(256) not null,
    provider_name varchar(100) not null,
    email varchar(256) not null,

    constraint ix_users_provider_data unique (provider_id, provider_name)
);

create table platforms (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
    name varchar(200) not null constraint ix_platform_name unique,
    short_name varchar(5) not null constraint ix_platform_short_name unique
);

create table games (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
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

create function check_user_playthrough(user_id uuid, playthrough_id uuid)
    returns boolean
    as $$
        begin
            return exists(
                select from playthroughs p
                join games g on p.game_id = g.id
                where p.id = check_user_playthrough.playthrough_id and g.user_id = check_user_playthrough.user_id
            );
        end;
    $$
    language plpgsql;
