create table game_notes (
    id uuid primary key default gen_random_uuid(),
    game_id uuid not null references games(id),
    title varchar(100) not null,
    value varchar(10000) not null,
    kind smallint not null, --text note, link, and so on...
    added_on timestamp with time zone not null default now(),
    pinned boolean not null
);

create function check_user_note(user_id uuid, note_id uuid)
    returns boolean
as $$
begin
    return exists(
        select from game_notes n
                        join games g on n.game_id = g.id
        where n.id = check_user_note.note_id and g.user_id = check_user_note.user_id
    );
end;
$$
    language plpgsql;
