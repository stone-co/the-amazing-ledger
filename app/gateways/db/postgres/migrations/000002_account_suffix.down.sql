BEGIN;

ALTER TABLE entries
    DROP COLUMN account_suffix
;

COMMIT;
