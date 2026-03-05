CREATE TABLE agendamentos (
    id VARCHAR(20) PRIMARY KEY,

    cliente_id   VARCHAR(20) NOT NULL,
    prestador_id VARCHAR(20) NOT NULL,
    catalogo_id  VARCHAR(20) NOT NULL,

    data_hora_inicio TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    data_hora_fim    TIMESTAMP WITHOUT TIME ZONE NOT NULL,

    status INTEGER NOT NULL,
    notas TEXT,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_agendamento_cliente
        FOREIGN KEY (cliente_id)
        REFERENCES clientes (id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_agendamento_prestador
        FOREIGN KEY (prestador_id)
        REFERENCES prestadores (id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_agendamento_catalogo
        FOREIGN KEY (catalogo_id)
        REFERENCES catalogos (id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_data_hora_valida
        CHECK (data_hora_inicio < data_hora_fim),

    CONSTRAINT chk_status_valido
        CHECK (status IN (1, 2, 3, 4))
);

CREATE INDEX idx_agendamentos_prestador_data
ON agendamentos (prestador_id, data_hora_inicio);

CREATE INDEX idx_agendamentos_cliente
ON agendamentos (cliente_id);

CREATE INDEX idx_agendamentos_status
ON agendamentos (status);