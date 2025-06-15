BEGIN;
DROP TRIGGER IF EXISTS update_users_updated_at ON auth.users;
DROP FUNCTION IF EXISTS auth.update_updated_at_column();
COMMIT;