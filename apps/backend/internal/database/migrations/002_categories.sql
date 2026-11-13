-- categories table
CREATE TABLE IF NOT EXISTS categories (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  name varchar(100) NOT NULL,
  user_id uuid NOT NULL,
  color char(7) NOT NULL DEFAULT '#6B7280',
  created_at timestamptz NOT NULL DEFAULT now (),
  updated_at timestamptz NOT NULL DEFAULT now ()
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_categories_created_at ON categories (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories (user_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_user_id_name ON categories (user_id, name);

-- trigger
CREATE OR REPLACE TRIGGER set_updated_at_categories BEFORE
UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();
