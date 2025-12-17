-- Drop indexes
DROP INDEX IF EXISTS idx_documents_file_type;
DROP INDEX IF EXISTS idx_documents_deleted_at;

-- Drop documents table
DROP TABLE IF EXISTS "documents";

-- Drop document_type enum
DROP TYPE IF EXISTS document_type;
