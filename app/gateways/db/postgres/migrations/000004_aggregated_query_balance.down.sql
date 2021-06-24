begin;

drop table if exists aggregated_query_balance;

drop function if exists get_aggregated_account_balance;

drop function if exists _get_aggregated_query_balance;
drop function if exists _get_aggregated_query_balance_since;

drop procedure if exists _update_aggregated_query_balance;
drop procedure if exists _insert_aggregated_query_balance;

commit;
