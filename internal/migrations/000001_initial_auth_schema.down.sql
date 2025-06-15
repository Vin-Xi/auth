BEGIN;
DROP SCHEMA IF EXISTS auth CASCADE;
-- We generally do not drop extensions in a down migration
-- as other parts of the system might depend on it.
COMMIT;