BEGIN;

DROP TABLE entries;
DROP TYPE  operation_type;
DROP TABLE accounts;
DROP TYPE  account_type;

DROP INDEX idx_entries_account_class;
DROP INDEX idx_entries_account_group;
DROP INDEX idx_entries_account_subgroup;
DROP INDEX idx_entries_account_id;

COMMIT;
