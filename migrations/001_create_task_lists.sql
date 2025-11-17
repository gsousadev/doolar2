-- Criar tabela task_lists
CREATE TABLE IF NOT EXISTS task_lists (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    tasks JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índice para busca por título
CREATE INDEX IF NOT EXISTS idx_task_lists_title ON task_lists(title);

-- Trigger para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_task_lists_updated_at 
    BEFORE UPDATE ON task_lists 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
