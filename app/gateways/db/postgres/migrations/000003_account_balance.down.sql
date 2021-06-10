begin;

drop table if exists account_balance;

drop function if exists get_account_balance(_account ltree, out _credit_balance bigint, out _debit_balance bigint);

drop function if exists _get_account_balance(_account ltree);
drop function if exists _get_account_balance_since(_account ltree, _dt timestamptz);

drop procedure if exists _update_account_balance;
drop procedure if exists _insert_account_balance;

commit;
