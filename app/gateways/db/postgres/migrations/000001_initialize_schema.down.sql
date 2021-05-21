begin;

drop index idx_entry_account_gist;
drop index idx_entry_company;
drop index idx_entry_event;
drop index idx_created_at;
drop index idx_competence_date;

drop table entry;
drop table event;
drop table company;

drop type operation;

drop extension ltree;

commit;
