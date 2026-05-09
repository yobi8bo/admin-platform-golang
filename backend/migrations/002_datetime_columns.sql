DO $$
DECLARE
  item RECORD;
BEGIN
  FOR item IN
    SELECT table_name, column_name
    FROM (VALUES
      ('sys_depts', 'created_at'),
      ('sys_depts', 'updated_at'),
      ('sys_users', 'created_at'),
      ('sys_users', 'updated_at'),
      ('sys_roles', 'created_at'),
      ('sys_roles', 'updated_at'),
      ('sys_menus', 'created_at'),
      ('sys_menus', 'updated_at'),
      ('sys_files', 'created_at'),
      ('sys_files', 'updated_at'),
      ('sys_login_logs', 'created_at'),
      ('sys_operation_logs', 'created_at')
    ) AS columns(table_name, column_name)
  LOOP
    IF EXISTS (
      SELECT 1
      FROM information_schema.columns
      WHERE table_schema = 'public'
        AND table_name = item.table_name
        AND column_name = item.column_name
        AND data_type = 'bigint'
    ) THEN
      EXECUTE format(
        'ALTER TABLE %I ALTER COLUMN %I DROP DEFAULT',
        item.table_name,
        item.column_name
      );
      EXECUTE format(
        'ALTER TABLE %I ALTER COLUMN %I TYPE TIMESTAMPTZ USING to_timestamp(%I / 1000.0)',
        item.table_name,
        item.column_name,
        item.column_name
      );
      EXECUTE format(
        'ALTER TABLE %I ALTER COLUMN %I SET DEFAULT CURRENT_TIMESTAMP',
        item.table_name,
        item.column_name
      );
    END IF;
  END LOOP;
END $$;
