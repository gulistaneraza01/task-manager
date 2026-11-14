-- status enum
CREATE TYPE todo_status AS ENUM ('draft', 'active', 'completed', 'archived');

-- priority enum
CREATE TYPE todo_priority AS ENUM ('low', 'medium', 'high');

-- todos table
CREATE TABLE IF NOT EXISTS todos (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  title varchar(255) NOT NULL,
  description text,
  status todo_status NOT NULL DEFAULT 'draft',
  priority todo_priority NOT NULL DEFAULT 'medium',
  due_date timestamptz,
  completed_at timestamptz,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  sort_order serial,
  category_id uuid references categories (id),
  parent_todo_id uuid references todos (id) ON DELETE CASCADE CONSTRAINT no_self_parent CHECK (
    parent_todo_id IS DISTINCT
    FROM
      id
  ),
  user_id uuid NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now (),
  updated_at timestamptz NOT NULL DEFAULT now ()
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_todos_category_id ON todos (category_id);

CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos (user_id);

CREATE INDEX IF NOT EXISTS idx_todos_parent_id ON todos (parent_todo_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_todos_status ON todos (status);

CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos (priority);

-- trigger
CREATE OR REPLACE TRIGGER set_updated_at_todos BEFORE
UPDATE ON todos FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();
