CREATE TABLE prestadores (
    id VARCHAR(20) PRIMARY KEY,
    nome VARCHAR(100) NOT NULL CHECK (char_length(nome) BETWEEN 3 AND 100),
    cpf CHAR(11) NOT NULL,
    email VARCHAR(255),
    telefone VARCHAR(15) NOT NULL CHECK (char_length(telefone) BETWEEN 8 AND 15),
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    imagem_url VARCHAR(200),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_prestadores_cpf
ON prestadores (cpf);

CREATE UNIQUE INDEX uq_prestadores_email
ON prestadores (email)
WHERE email IS NOT NULL;

CREATE TABLE prestador_catalogos (
    prestador_id VARCHAR(20) NOT NULL,
    catalogo_id  VARCHAR(20) NOT NULL,

    PRIMARY KEY (prestador_id, catalogo_id),

    CONSTRAINT fk_prestador
        FOREIGN KEY (prestador_id)
        REFERENCES prestadores (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_catalogo
        FOREIGN KEY (catalogo_id)
        REFERENCES catalogos (id)
        ON DELETE RESTRICT
);
