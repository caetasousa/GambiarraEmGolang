CREATE TABLE clientes (
    id         VARCHAR(20) PRIMARY KEY,
    nome       VARCHAR(100) NOT NULL CHECK (char_length(nome) BETWEEN 3 AND 100),
    email      VARCHAR(255) NOT NULL,
    telefone   VARCHAR(15)  NOT NULL CHECK (char_length(telefone) BETWEEN 8 AND 15),
    senha_hash VARCHAR(255) NOT NULL DEFAULT '',
    role       VARCHAR(20)  NOT NULL DEFAULT 'cliente',
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_cliente_admin_role CHECK (role IN ('cliente', 'admin'))
);

CREATE UNIQUE INDEX uq_clientes_email
ON clientes (email);

-- Promove o usuário especificado para admin se ele existir (útil para desenvolvimento)
COMMENT ON COLUMN clientes.senha_hash IS 'Senha hasheada com bcrypt';
COMMENT ON COLUMN clientes.role IS 'Role do usuário - cliente ou admin';

