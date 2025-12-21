CREATE TABLE catalogos (
    id VARCHAR(20) PRIMARY KEY,
    nome VARCHAR(100) NOT NULL CHECK (char_length(nome) BETWEEN 3 AND 100),
    duracao_padrao INTEGER NOT NULL CHECK (duracao_padrao > 0),
    preco INTEGER NOT NULL CHECK (preco >= 0),
    categoria VARCHAR(50) NOT NULL CHECK (char_length(categoria) BETWEEN 3 AND 50),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);