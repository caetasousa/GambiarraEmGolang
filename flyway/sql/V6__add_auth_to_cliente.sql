-- V6__add_auth_to_cliente.sql
-- Adiciona campos de autenticação na tabela clientes

-- Adicionar coluna de senha (hash bcrypt)
ALTER TABLE clientes
ADD COLUMN senha_hash VARCHAR(255) NOT NULL DEFAULT '';

-- Adicionar coluna de role (sempre 'cliente' para clientes)
ALTER TABLE clientes
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'cliente';

-- Adicionar constraint para garantir que role seja sempre 'cliente'
ALTER TABLE clientes
ADD CONSTRAINT chk_cliente_role
CHECK (role = 'cliente');

-- Comentários para documentação
COMMENT ON COLUMN clientes.senha_hash IS 'Senha hasheada com bcrypt';
COMMENT ON COLUMN clientes.role IS 'Role do usuário - sempre cliente';
