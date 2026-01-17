-- V7__add_auth_to_prestador.sql
-- Adiciona campos de autenticação na tabela prestadores

-- Adicionar coluna de senha (hash bcrypt)
ALTER TABLE prestadores
ADD COLUMN senha_hash VARCHAR(255) NOT NULL DEFAULT '';

-- Adicionar coluna de role (sempre 'prestador' para prestadores)
ALTER TABLE prestadores
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'prestador';

-- Adicionar constraint para garantir que role seja sempre 'prestador'
ALTER TABLE prestadores
ADD CONSTRAINT chk_prestador_role
CHECK (role = 'prestador');

-- Comentários para documentação
COMMENT ON COLUMN prestadores.senha_hash IS 'Senha hasheada com bcrypt';
COMMENT ON COLUMN prestadores.role IS 'Role do usuário - sempre prestador';
