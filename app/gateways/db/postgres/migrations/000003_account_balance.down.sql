begin;

drop table if exists account_balance;

drop function if exists get_account_balance;

drop function if exists _get_account_balance;
drop function if exists _get_account_balance_since;

drop procedure if exists _update_account_balance;
drop procedure if exists _insert_account_balance;

commit;
