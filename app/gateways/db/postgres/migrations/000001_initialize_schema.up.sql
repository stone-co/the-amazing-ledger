begin;

create extension if not exists ltree;

create table if not exists event
(
    id   smallint primary key generated always as identity,
    name text not null unique
);

create table if not exists entry
(
    id              uuid primary key,
    tx_id           uuid        not null,
    event           smallint    not null references event(id),
    operation       smallint    not null check (operation = 1 or operation = 2),
    version         int         not null,
    amount          bigint      not null,
    created_at      timestamptz not null default now(),
    competence_date timestamptz not null,
    account         ltree       not null,
    company         text        not null,
    metadata        jsonb       not null default '{}'
);

create index if not exists idx_entry_account_gist
    on entry using gist (account gist_ltree_ops(siglen=32));
create index if not exists idx_entry_tx
    on entry using btree (tx_id);
create index if not exists idx_entry_company
    on entry using btree (company);
create index if not exists idx_entry_event
    on entry using btree (event);
create index if not exists idx_created_at
    on entry using brin (created_at) with (pages_per_range = 32);

create table if not exists account_version (
    version int,
    account ltree primary key
) with (fillfactor = 50);

create or replace function update_account_version()
    returns trigger
    language plpgsql
as
$$
begin
    if new.version = 0 then
        update account_version set version = version + 1 where account = new.account returning version into new.version;
    else
        update account_version set version = new.version where account = new.account;
    end if;

    if not found then
        insert into account_version(version, account) values (1, new.account);
        new.version = 1;
    end if;

    return new;
end;
$$;

drop trigger if exists tg_update_account_version on entry;
create trigger tg_update_account_version
    before insert
    on entry
    for each row
    when (new.version >= 0)
execute procedure update_account_version();

create or replace function invalid_account_version()
    returns trigger
    language plpgsql
as
$$
begin
    raise exception 'invalid account version (from % to %)', old.version, new.version;
end;
$$;

create trigger tg_check_account_version
    before update
    on account_version
    for each row
    when (new.version != old.version + 1)
execute procedure invalid_account_version();

commit;
