CREATE TABLE IF NOT EXISTS documents (
                                         id BIGSERIAL PRIMARY KEY,
                                         external_id UUID NOT NULL DEFAULT gen_random_uuid(),
                                         user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                         document_type VARCHAR(50) NOT NULL,
                                         filename VARCHAR(255) NOT NULL,
                                         content_type VARCHAR(100) NOT NULL,
                                         file_content TEXT NOT NULL, -- Base64 do arquivo
                                         categories JSONB DEFAULT '[]', -- Categorias do documento em formato JSON
                                         status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, processing, processed, failed
                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_documents_external_id ON documents(external_id);
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_status ON documents(status);
CREATE INDEX idx_documents_document_type ON documents(document_type);