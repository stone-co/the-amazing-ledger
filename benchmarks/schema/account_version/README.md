When deciding how to handle transactional guarantees (when desirable) for each account, a possible solution would be
to have a table with only two columns: account and version.

This table would need to handle a lot of inserts at first, but as time goes on, it is expected that the number of
updates become higher than the number of inserts, so this table should be designed for handling numerous updates.

One of the solutions that postgres provides us is using [**HOT**](https://www.cybertec-postgresql.com/en/hot-updates-in-postgresql-for-better-performance/)
updates. Here we try to do a synthetic benchmark to see how (and if) this strategy suits our needs.

Our table should look something like this:

```sql
create table account_version (
    account ltree primary key,
    version int
);
```

With a little of [column tetris](https://www.2ndquadrant.com/en/blog/on-rocks-and-sand/) we can see that reordering
the columns gives us free 4 extra bytes for each entry (yay).

To use **HOT** updates, we should set a `fillfactor` lower than 100 (postgres default). In the end, we have this schema:

```sql
create table account_version (
    version int,
    account ltree primary key
) with (fillfactor = 50);
```

To run our benchmarks and see if this strategy has any effect for our use case we ran the tests below:

**Step 1.**

Create fake test data with 10 million accounts, divided between 4 categories:

- Really high frequency used accounts -> 1000 accounts
- High frequency used accounts -> 10%
- Medium frequency used accounts -> 30%
- Low frequency used accounts -> 60%

```sql
--
-- load ltree extension if needed
--

create extension if not exists ltree;

--
-- helper function to generate random bounded integers
--

CREATE OR REPLACE FUNCTION random_between(low INT, high INT)
    RETURNS INT AS
$$
BEGIN
    RETURN floor(random() * (high - low + 1) + low);
END;
$$ language 'plpgsql' STRICT;

--
-- create account frequency tables
--

create table test_accounts_really_high_frequency
(
    account ltree primary key
);

create table test_accounts_high_frequency
(
    account ltree primary key
);

create table test_accounts_medium_frequency
(
    account ltree primary key
);

create table test_accounts_low_frequency
(
    account ltree primary key
);

--
-- Generated 10 million random accounts
--

do
$$
    declare
        id         ltree;
        iterations integer := 10000000;
        counter    integer := 0;
        random     integer := 0;
    begin
        loop
            exit when counter = iterations;

            id := cast(replace(cast(gen_random_uuid() as text), '-', '_') as ltree);
            random := random_between(0, 10);

            if (random <= 6) then
                insert into test_accounts_low_frequency(account) values (id);
            elsif (random > 6 and random <= 9) then
                insert into test_accounts_medium_frequency(account) values (id);
            else
                insert into test_accounts_high_frequency(account) values (id);
            end if;

            counter := counter + 1;
        end loop;
    end;
$$;

--
-- select 1000 accounts to be frequently updated
--

insert into test_accounts_really_high_frequency (account)
select account
from test_accounts_high_frequency
limit 1000;

delete
from test_accounts_high_frequency
where account in (select account from test_accounts_really_high_frequency);
```

**Step 2.**

Create and fill tables that we want to benchmark:

```sql

--
-- default settings table
--

create table account_version (
    version int,
    account ltree primary key
);

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from account_version;

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from account_version;

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from account_version;

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from account_version;

--
-- custom fillfactor table
--

create table account_version_fillfactor (
    version int,
    account ltree primary key
) with (fillfactor = 50);

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from test_accounts_low_frequency;

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from test_accounts_medium_frequency;

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from test_accounts_high_frequency;

insert into account_version_fillfactor(version, account)
select random_between(1, 10), account
from test_accounts_really_high_frequency;
```

**Step 3.**

Benchmark select operation with *pgbench*:

- Default
```sql
select version from account_version where account = (
    select account
    from test_accounts_low_frequency
    tablesample bernoulli(1)
    limit 1
);
```

- Custom fillfactor
```sql
select version from account_version_fillfactor where account = (
    select account
    from test_accounts_low_frequency
    tablesample bernoulli(1)
    limit 1
);
```

**Step 4.**

Benchmark update operation with *pgbench*:

- Default
```sql
do
$$
    declare
        acc ltree;
        ver integer;
        random integer;
    begin
        select random_between(0, 10) into random;

        if random <= 4 then
            select account
            into acc
            from test_accounts_really_high_frequency
            tablesample bernoulli(2)
            limit 1;
        elsif random <= 7 then
            select account
            into acc
            from test_accounts_high_frequency
            tablesample bernoulli(2)
            limit 1;
        elsif random <= 8 then
            select account
            into acc
            from test_accounts_medium_frequency
            tablesample bernoulli(2)
            limit 1;
        else
            select account
            into acc
            from test_accounts_low_frequency
            tablesample bernoulli(2)
            limit 1;
        end if;

        select version
        into ver
        from account_version
        where account = acc;

        insert into account_version (version, account) values (ver, acc)
        on conflict on constraint account_version_pkey do update
            set version = ver
        where account_version.account = acc;
    end;
$$;
```

- Custom fillfactor
```sql
do
$$
    declare
        acc ltree;
        ver integer;
        random integer;
    begin
        select random_between(0, 10) into random;

        if random <= 4 then
            select account
            into acc
            from test_accounts_really_high_frequency
            tablesample bernoulli(2)
            limit 1;
        elsif random <= 7 then
            select account
            into acc
            from test_accounts_high_frequency
            tablesample bernoulli(2)
            limit 1;
        elsif random <= 8 then
            select account
            into acc
            from test_accounts_medium_frequency
            tablesample bernoulli(2)
            limit 1;
        else
            select account
            into acc
            from test_accounts_low_frequency
            tablesample bernoulli(2)
            limit 1;
        end if;

        select version
        into ver
        from account_version_fillfactor
        where account = acc;

        insert into account_version_fillfactor (version, account) values (ver, acc)
        on conflict on constraint account_version_fillfactor_pkey do update
            set version = ver
        where account_version_fillfactor.account = acc;
    end;
$$;
```

**Results**:
--

Running the benchmark above with the following configurations:

- Machine:
    - CPU: Ryzen 3700x
    - RAM: 16GB DDR4 2400MHz
    - Storage: 128GB SATA SSD
    
- Postgres:
    - version: 13.2
    - max_connections = '100'
    - shared_buffers = '4GB'
    - effective_cache_size = '12GB'
    - maintenance_work_mem = '1GB'
    - checkpoint_completion_target = '0.9'
    - wal_buffers = '16MB'
    - default_statistics_target = '100'
    - random_page_cost = '1.1'
    - effective_io_concurrency = '200'
    - work_mem = '20971kB'
    - min_wal_size = '2GB'
    - max_wal_size = '8GB'
    - max_worker_processes = '8'
    - max_parallel_workers_per_gather = '4'
    - max_parallel_workers = '8'
    - max_parallel_maintenance_workers = '4'
    
- PGBench:
    - 95 clients (-c)
    - 12 threads (-j)
    - select bench for 10 minutes
    - update bench for 15 minutes


- *Select results*:

|                    | latency average | latency stddev | tps (including connection establishing) | tps (excluding connection establishing) |
|--------------------|-----------------|----------------|-----------------------------------------|-----------------------------------------|
| default            | 0.830 ms        | 0.717 ms       | 111622.296396                           | 111623.350941                           |
| custom fill factor | 0.829 ms        | 0.715 ms       | 111825.402327                           | 111827.162195                           |


- *Update results*:

|                    | latency average | latency stddev | tps (including connection establishing) | tps (excluding connection establishing) |
|--------------------|-----------------|----------------|-----------------------------------------|-----------------------------------------|
| default            | 4.763 ms        | 3.239 ms       | 19926.127432                            | 19926.369692                            |
| custom fill factor | 4.458 ms        | 2.481 ms       | 21267.856228                            | 21268.063340                            |

We can take a look at table statistics, to see why the custom fill factor option performed better, with the query below:

```sql
select 
    relname,
    n_tup_upd,
    n_tup_hot_upd,
    n_tup_upd - n_tup_hot_upd delta,
    n_dead_tup
from pg_stat_all_tables
where relname in ('account_version', 'account_version_fillfactor');
```

Here the `n_dead_tup` represents how many dead tuples the table have. This number was consistently lower
though the test runs for the custom fillfactor table and `delta` showed the same behaviour. This tell us that **HOT**
worked for our use case, with less dead tuples by the end, which leads to a lower maintenance
overhead by the vacuuming process.

**DISCLAIMER**

The results above are an example of one of our outcome, although throughout our tests runs the custom fill factor
option was always the winner.
