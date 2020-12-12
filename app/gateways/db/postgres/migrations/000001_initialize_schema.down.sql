BEGIN;

DROP INDEX idx_entries_account_class;
DROP INDEX idx_entries_account_group;
DROP INDEX idx_entries_account_subgroup;
DROP INDEX idx_entries_account_id;

DROP TABLE entries;
DROP TYPE  operation_type;
DROP TYPE  account_class;

COMMIT;
