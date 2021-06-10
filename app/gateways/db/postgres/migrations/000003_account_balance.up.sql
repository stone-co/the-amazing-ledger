begin;

create table if not exists account_balance
(
    credit_balance bigint      not null,
    debit_balance  bigint      not null,
    tx_date        timestamptz not null,
    account        ltree primary key
) with (fillfactor = 50);

create or replace function _get_account_balance(_account ltree)
    returns table
        (
            credit bigint,
            debit  bigint,
            tx_dt  timestamptz
        )
    language sql
as
$$
select
    sub.total_debit,
    sub.total_credit,
    sub.recent
from (
        select
            sum(
                case operation
                    when 1 then amount
                    else 0
                end
            ) as total_credit,
            sum(
                case operation
                    when 2 then amount
                    else 0
                end
            ) as total_debit,
            (max(array[created_at]))[1] as recent
        from
            entry
        where
            account = _account
    ) sub;
$$ stable rows 1
;

create or replace function _get_account_balance_since(_account ltree, _dt timestamptz)
    returns table
        (
            credit bigint,
            debit  bigint,
            tx_dt  timestamptz
        )
    language sql
as
$$
select
    sub.total_debit,
    sub.total_credit,
    sub.recent
from (
    select
        sum(
            case operation
                when 1 then amount
                else 0
            end
        ) as total_credit,
        sum(
            case operation
                when 2 then amount
                else 0
            end
        ) as total_debit,
        (max(array[created_at]))[1] as recent
    from
        entry
    where
        account = _account
        and created_at > _dt
    ) sub;
$$ stable rows 1
;

create or replace procedure _update_account_balance(
    _account ltree, _credit bigint, _debit bigint, _dt timestamptz
)
    language sql
as
$$
update account_balance
set
    credit_balance = _credit,
    debit_balance = _debit,
    tx_date = _dt
where
    account = _account;
$$;

create or replace procedure _insert_account_balance(
    _account ltree, _credit bigint, _debit bigint, _dt timestamptz
)
    language sql
as
$$
insert into account_balance (credit_balance, debit_balance, tx_date, account)
values (_credit, _debit, _dt, _account);
$$;

create or replace function get_account_balance(_account ltree, out _credit_balance bigint, out _debit_balance bigint)
    returns record
    language plpgsql
as
$$
declare
    _current_credit  bigint;
    _current_debit   bigint;
    _last_tx_date    timestamptz;
begin
    select
        credit_balance,
        debit_balance,
        tx_date
    into
        _current_credit,
        _current_debit,
        _last_tx_date
    from
         account_balance
    where
        account = _account;

    if (_last_tx_date is null) then
        select
            credit,
            debit,
            tx_dt
        into
            _credit_balance,
            _debit_balance,
            _last_tx_date
        from
            _get_account_balance(_account);

        call _insert_account_balance(
            _account => _account,
            _credit => _credit_balance,
            _debit => _debit_balance,
            _dt => _last_tx_date
        );

        return;
    end if;

    select
        credit + _current_credit,
        debit + _current_debit,
        tx_dt
    into
        _credit_balance,
        _debit_balance,
        _last_tx_date
    from
        _get_account_balance_since(_account, _last_tx_date);

    if (_credit_balance is null) then
        select
            _current_credit,
            _current_debit
        into
            _credit_balance,
            _debit_balance;

        return;
    end if;

    call _update_account_balance(
        _account => _account,
        _credit => _credit_balance,
        _debit => _debit_balance,
        _dt => _last_tx_date
    );
end;
$$ volatile;

commit;
