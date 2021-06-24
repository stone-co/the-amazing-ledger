begin;

create table if not exists account_balance
(
    credit      bigint      not null,
    debit       bigint      not null,
    tx_date     timestamptz not null,
    account     ltree primary key
) with (fillfactor = 70);

create or replace function _get_account_balance(_account ltree)
    returns table
        (
            partial_credit bigint,
            partial_debit  bigint,
            partial_date   timestamptz,
            recent_credit  bigint,
            recent_debit   bigint,
            recent_version int
        )
    language sql
as
$$
    select
        sum(sub.total_credit) filter (where sub.row_number > 1) as partial_credit,
        sum(sub.total_debit)  filter (where sub.row_number > 1) as partial_debit,
        max(created_at)       filter (where sub.row_number = 2) as partial_date,

        sum(sub.total_credit) filter (where sub.row_number = 1) as recent_credit,
        sum(sub.total_debit)  filter (where sub.row_number = 1) as recent_debit,
        max(version)          filter (where sub.row_number = 1) as recent_version
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
            max(version) as version,
            created_at,
            row_number() over (order by created_at desc) as row_number
        from entry
        where account = _account
        group by created_at
        order by created_at desc
    ) sub
$$ stable rows 1
;

create or replace function _get_account_balance_since(_account ltree, _dt timestamptz)
    returns table
        (
            partial_credit bigint,
            partial_debit  bigint,
            partial_date   timestamptz,
            recent_credit  bigint,
            recent_debit   bigint,
            recent_version int
        )
    language sql
as
$$
    select
        sum(sub.total_credit) filter (where sub.row_number > 1) as partial_credit,
        sum(sub.total_debit)  filter (where sub.row_number > 1) as partial_debit,
        max(created_at)       filter (where sub.row_number = 2) as partial_date,

        sum(sub.total_credit) filter (where sub.row_number = 1) as recent_credit,
        sum(sub.total_debit)  filter (where sub.row_number = 1) as recent_debit,
        max(version)          filter (where sub.row_number = 1) as recent_version
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
            max(version) as version,
            created_at,
            row_number() over (order by created_at desc) as row_number
        from entry
        where
            account = _account
            and created_at > _dt
        group by created_at
        order by created_at desc
    ) sub
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
        credit = _credit,
        debit = _debit,
        tx_date = _dt
    where account = _account;
$$;

create or replace procedure _insert_account_balance(
    _account ltree, _credit bigint, _debit bigint, _dt timestamptz
)
    language sql
as
$$
    insert into account_balance (credit, debit, tx_date, account)
    values (_credit, _debit, _dt, _account);
$$;

create or replace function get_account_balance(
    in _account ltree,
    out credit_balance bigint, out debit_balance bigint, out version int
)
    returns record
    language plpgsql
as
$$
declare
    _existing_credit bigint;
    _existing_debit  bigint;
    _existing_date   timestamptz;

    _partial_credit bigint;
    _partial_debit  bigint;
    _partial_date   timestamptz;
begin
    select
        credit,
        debit,
        tx_date
    into
        _existing_credit,
        _existing_debit,
        _existing_date
    from
        account_balance
    where
        account = _account;

    if (_existing_credit is null) then
        select
            partial_credit,
            partial_debit,
            partial_date,

            coalesce(partial_credit, 0) + recent_credit,
            coalesce(partial_debit, 0) + recent_debit,
            recent_version
        into
            _partial_credit,
            _partial_debit,
            _partial_date,

            credit_balance,
            debit_balance,
            version
        from
            _get_account_balance(_account);

        if (version is null) then
            raise no_data_found;
        elsif (_partial_credit is null) then
            return;
        end if;

        call _insert_account_balance(
            _account => _account,
            _credit => _partial_credit,
            _debit => _partial_debit,
            _dt => _partial_date
        );

        return;
    end if;

    select
        _existing_credit + partial_credit,
        _existing_debit + partial_debit,
        partial_date,

        _existing_credit + coalesce(partial_credit, 0) + coalesce(recent_credit, 0),
        _existing_debit + coalesce(partial_debit, 0) + coalesce(recent_debit, 0),
        recent_version
    into
        _partial_credit,
        _partial_debit,
        _partial_date,

        credit_balance,
        debit_balance,
        version
    from
        _get_account_balance_since(_account, _existing_date);

    if (_partial_date is null) then
        return;
    end if;

    call _update_account_balance(
        _account => _account,
        _credit => _partial_credit,
        _debit => _partial_debit,
        _dt => _partial_date
    );
end;
$$ volatile;

commit;
