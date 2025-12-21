CREATE TABLE clientes (
    id         VARCHAR(20) PRIMARY KEY,
    nome       VARCHAR(100) NOT NULL CHECK (char_length(nome) BETWEEN 3 AND 100),
    email      VARCHAR(255) NOT NULL,
    telefone   VARCHAR(15)  NOT NULL CHECK (char_length(telefone) BETWEEN 8 AND 15),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_clientes_email
ON clientes (email);

