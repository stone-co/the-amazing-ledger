begin;

create table if not exists account_balance
(
    credit      bigint      not null,
    debit       bigint      not null,
    tx_date     timestamptz not null,
    tx_version  int         not null,
    account     ltree primary key
) with (fillfactor = 50);

create or replace function _get_account_balance(_account ltree)
    returns table
        (
            credit      bigint,
            debit       bigint,
            tx_dt       timestamptz,
            tx_version  int,
            min_version int
        )
    language sql
as
$$
    select
        sub.total_debit,
        sub.total_credit,
        to_timestamp(sub.recent[1]) as last_tx_date,
        sub.recent[2]::int as last_tx_version,
        sub.min_version
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
            (max(array[extract('epoch' from created_at), version])) as recent,
            (min(array[version]))[1] as min_version
        from entry
        where account = _account
    ) sub;
$$ stable rows 1
;

create or replace function _get_account_balance_since(_account ltree, _dt timestamptz, _version int)
    returns table
        (
            credit      bigint,
            debit       bigint,
            tx_dt       timestamptz,
            tx_version  int,
            min_version int
        )
    language sql
as
$$
    select
    sub.total_debit,
    sub.total_credit,
    to_timestamp(sub.recent[1]) as last_tx_date,
    sub.recent[2]::int as last_tx_version,
    sub.min_version
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
            (max(array[extract('epoch' from created_at), version])) as recent,
            (min(array[version]))[1] as min_version
        from entry
        where
            account = _account
            and created_at >= _dt
            and (created_at, version) > (_dt, _version)
    ) sub;
$$ stable rows 1
;

create or replace procedure _update_account_balance(
    _account ltree, _credit bigint, _debit bigint, _dt timestamptz, _version int
)
    language sql
as
$$
    update account_balance
    set
        credit = _credit,
        debit = _debit,
        tx_date = _dt,
        tx_version = _version
    where account = _account;
$$;

create or replace procedure _insert_account_balance(
    _account ltree, _credit bigint, _debit bigint, _dt timestamptz, _version int
)
    language sql
as
$$
    insert into account_balance (credit, debit, tx_date, tx_version, account)
    values (_credit, _debit, _dt, _version, _account);
$$;

create or replace function get_account_balance(
    in _account ltree,
    out credit_balance bigint, out debit_balance bigint, out dt timestamptz, out version int
)
    returns record
    language plpgsql
as
$$
declare
    _current_credit  bigint;
    _current_debit   bigint;
    _last_tx_date    timestamptz;
    _last_tx_version int;
    _min_tx_version  int;
begin
    select
        credit,
        debit,
        tx_date,
        tx_version
    into
        _current_credit,
        _current_debit,
        _last_tx_date,
        _last_tx_version
    from
        account_balance
    where
        account = _account;

    if (_last_tx_version is null) then
        select
            credit,
            debit,
            tx_dt,
            tx_version,
            min_version
        into
            credit_balance,
            debit_balance,
            dt,
            version,
            _min_tx_version
        from
            _get_account_balance(_account);

        if (_min_tx_version <= 0) then
            return;
        end if;

        call _insert_account_balance(
            _account => _account,
            _credit => credit_balance,
            _debit => debit_balance,
            _dt => dt,
            _version => version
        );

        return;
    end if;

    select
        coalesce(credit + _current_credit, _current_credit),
        coalesce(debit + _current_debit, _current_debit),
        coalesce(tx_dt, _last_tx_date),
        coalesce(tx_version, _last_tx_version),
        min_version
    into
        credit_balance,
        debit_balance,
        dt,
        version,
        _min_tx_version
    from
        _get_account_balance_since(_account, _last_tx_date, _last_tx_version);

    if (version is null or _min_tx_version <= 0) then
        return;
    end if;

    call _update_account_balance(
        _account => _account,
        _credit => credit_balance,
        _debit => debit_balance,
        _dt => dt,
        _version => version
    );
end;
$$ volatile;

commit;
