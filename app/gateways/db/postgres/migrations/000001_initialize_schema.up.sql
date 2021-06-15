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

commit;
