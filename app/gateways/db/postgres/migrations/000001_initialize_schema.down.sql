begin;

drop index if exists idx_entry_account_gist;
drop index if exists idx_entry_tx;
drop index if exists idx_entry_company;
drop index if exists idx_entry_event;
drop index if exists idx_created_at;

drop table if exists entry;
drop table if exists event;

drop table if exists account_version_fillfactor;

drop extension if exists ltree;

commit;
