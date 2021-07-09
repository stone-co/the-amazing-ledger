begin;

create table if not exists aggregated_query_balance
(
    balance bigint      not null,
    tx_date timestamptz not null,
    query   text primary key
) with (fillfactor = 70);

create or replace function _get_aggregated_query_balance(_query lquery)
    returns table
            (
                partial_balance bigint,
                partial_date    timestamptz,
                recent_balance  bigint
            )
    language sql
as
$$
select sum(sub.balance) filter (where sub.row_number > 1) as partial_balance,
       max(created_at) filter (where sub.row_number = 2)  as partial_date,
       sum(sub.balance) filter (where sub.row_number = 1) as recent_balance
from (
         select coalesce(sum(amount) filter (where operation = 1), 0) -
                coalesce(sum(amount) filter (where operation = 2), 0) as balance,
                created_at,
                row_number() over (order by created_at desc)          as row_number
         from entry
         where account ~ _query
         group by created_at
         order by created_at desc
     ) sub
$$ stable
   rows 1
;

create or replace function _get_aggregated_query_balance_since(_query lquery, _dt timestamptz)
    returns table
            (
                partial_balance bigint,
                partial_date    timestamptz,
                recent_balance  bigint
            )
    language sql
as
$$
select sum(sub.balance) filter (where sub.row_number > 1) as partial_balance,
       max(created_at) filter (where sub.row_number = 2)  as partial_date,
       sum(sub.balance) filter (where sub.row_number = 1) as recent_credit
from (
         select coalesce(sum(amount) filter (where operation = 1), 0) -
                coalesce(sum(amount) filter (where operation = 2), 0) as balance,
                created_at,
                row_number() over (order by created_at desc)          as row_number
         from entry
         where account ~ _query
           and created_at > _dt
         group by created_at
         order by created_at desc
     ) sub
$$ stable
   rows 1
;

create or replace procedure _update_aggregated_query_balance(
    _query text, _balance bigint, _dt timestamptz
)
    language sql
as
$$
update aggregated_query_balance
set balance = _balance,
    tx_date = _dt
where query = _query;
$$;

create or replace procedure _insert_aggregated_query_balance(
    _query text, _balance bigint, _dt timestamptz
)
    language sql
as
$$
insert into aggregated_query_balance (balance, tx_date, query)
values (_balance, _dt, _query);
$$;

create or replace function query_aggregated_account_balance(
    in _query lquery, out total_balance bigint
)
    returns bigint
    language plpgsql
as
$$
declare
    _existing_balance bigint;
    _existing_date    timestamptz;
    _partial_balance  bigint;
    _partial_date     timestamptz;
begin
    select balance,
           tx_date
    into
        _existing_balance,
        _existing_date
    from aggregated_query_balance
    where query = _query::text;

    if (_existing_balance is null) then
        select partial_balance,
               partial_date,
               coalesce(partial_balance, 0) + recent_balance
        into
            _partial_balance,
            _partial_date,
            total_balance
        from
            _get_aggregated_query_balance(_query);

        if (total_balance is null) then
            raise no_data_found;
        elsif (_partial_balance is null) then
            return;
        end if;

        call _insert_aggregated_query_balance(
                _query => _query::text,
                _balance => _partial_balance,
                _dt => _partial_date
            );

        return;
    end if;

    select _existing_balance + partial_balance,
           partial_date,
           _existing_balance + coalesce(partial_balance, 0) + coalesce(recent_balance, 0)
    into
        _partial_balance,
        _partial_date,
        total_balance
    from
        _get_aggregated_query_balance_since(_query, _existing_date);

    if (_partial_date is null) then
        return;
    end if;

    call _update_aggregated_query_balance(
            _query => _query::text,
            _balance => _partial_balance,
            _dt => _partial_date
        );
end;
$$ volatile;

commit;
