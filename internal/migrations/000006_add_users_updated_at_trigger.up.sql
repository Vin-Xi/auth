BEGIN;
CREATE OR REPLACE FUNCTION auth.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = (now() AT TIME ZONE 'utc');
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON auth.users
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();
COMMIT;