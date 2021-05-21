begin;

create extension if not exists ltree;

create table event
(
    id   smallint primary key generated always as identity,
    name text not null unique
);

create table company
(
    id   smallint primary key generated always as identity,
    name text not null unique
);

create type operation as enum ('debit', 'credit');

create table entry
(
    id              uuid primary key,
    tx_id           uuid        not null,
    version         bigint      not null,
    operation       operation   not null,
    company         smallint    not null references company(id),
    event           smallint    not null references event(id),
    amount          bigint      not null,
    created_at      timestamptz not null default now(),
    competence_date timestamptz not null,
    account         ltree       not null,
    metadata        jsonb       not null default '{}'
);

create index idx_entry_account_gist
    on entry using gist (account gist_ltree_ops(siglen=32));
create index idx_entry_company
    on entry using btree (company);
create index idx_entry_event
    on entry using btree (event);
create index idx_created_at
    on entry using brin (created_at);
create index idx_competence_date
    on entry using btree (competence_date);

create unique index idx_entry_account_max_version
    on entry (account, version desc);

commit;
