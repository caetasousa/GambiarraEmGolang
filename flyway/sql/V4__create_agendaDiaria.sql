CREATE TABLE agendas_diarias (
    id VARCHAR(20) PRIMARY KEY,
    prestador_id VARCHAR(20) NOT NULL,
    data DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_prestador
        FOREIGN KEY(prestador_id) REFERENCES prestadores(id)
        ON DELETE CASCADE,
    CONSTRAINT uq_prestador_data UNIQUE(prestador_id, data)
);

CREATE TABLE intervalos_diarios (
    id VARCHAR(20) PRIMARY KEY,
    agenda_id VARCHAR(20) NOT NULL,
    hora_inicio TIME NOT NULL,
    hora_fim TIME NOT NULL,
    CONSTRAINT fk_agenda
        FOREIGN KEY(agenda_id) REFERENCES agendas_diarias(id)
        ON DELETE CASCADE,
    CONSTRAINT chk_horario_valido CHECK (hora_inicio < hora_fim)
);