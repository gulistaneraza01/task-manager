CREATE TABLE IF NOT EXISTS comments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  todo_id uuid NOT NULL REFERENCES todos (id) ON DELETE CASCADE,
  user_id uuid NOT NULL,
  content text NOT NULL,
  parent_comment_id uuid REFERENCES comments (id) ON DELETE CASCADE CONSTRAINT no_self_reference CHECK (
    parent_comment_id IS DISTINCT
    FROM
      id
  ),
  created_at timestamptz NOT NULL DEFAULT now (),
  updated_at timestamptz NOT NULL DEFAULT now ()
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments (parent_comment_id);

CREATE INDEX IF NOT EXISTS idx_comments_todo_id_created_at ON comments (todo_id, created_at DESC);

-- trigger
CREATE OR REPLACE TRIGGER set_updated_at_comments BEFORE
UPDATE ON comments FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();
