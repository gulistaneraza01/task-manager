--document_stores table
CREATE TABLE IF NOT EXISTS document_stores (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id uuid NOT NULL,
  bucket varchar(255) NOT NULL,
  key varchar(255) NOT NULL,
  url text NOT NULL UNIQUE,
  file_name text NOT NULL,
  file_size bigint NOT NULL check (file_size > 0),
  file_type text NOT NULL,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now (),
  updated_at timestamptz NOT NULL DEFAULT now (),
  CONSTRAINT uq_document_stores_bucket_key UNIQUE (bucket, key)
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_document_stores_user_id ON document_stores (user_id);

CREATE INDEX IF NOT EXISTS idx_document_stores_bucket ON document_stores (bucket);

CREATE INDEX IF NOT EXISTS idx_document_stores_key ON document_stores (key);

CREATE INDEX IF NOT EXISTS idx_document_stores_url ON document_stores (url);
