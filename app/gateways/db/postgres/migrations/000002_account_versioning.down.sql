begin;

drop trigger if exists tg_update_account_version on entry;

drop table if exists account_version;

drop function if exists update_account_version;

commit;
