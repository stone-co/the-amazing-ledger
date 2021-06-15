begin;

create table if not exists account_version (
    version int,
    account ltree primary key
) with (fillfactor = 70);

create or replace function update_account_version()
    returns trigger
    language plpgsql
as
$$
begin
    if new.version = 0 then
        update account_version set version = version + 1 where account = new.account returning version into new.version;
    else
        update account_version set version = new.version where account = new.account;
    end if;

    if not found then
        insert into account_version(version, account) values (1, new.account);
        new.version = 1;
    end if;

    return new;
end;
$$;

create trigger tg_update_account_version
    before insert
    on entry
    for each row
    when (new.version >= 0)
execute procedure update_account_version();

create or replace function invalid_account_version()
    returns trigger
    language plpgsql
as
$$
begin
    raise exception 'invalid account version (from % to %)', old.version, new.version;
end;
$$;

create trigger tg_check_account_version
    before update
    on account_version
    for each row
    when (new.version != old.version + 1)
execute procedure invalid_account_version();

commit;
