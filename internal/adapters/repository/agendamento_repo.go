package repository

import (
	"database/sql"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendamentoPostgresRepository struct {
	db *sql.DB
}

func NovoAgendamentoPostgresRepository(db *sql.DB) port.AgendamentoRepositorio {
	return &AgendamentoPostgresRepository{db: db}
}

func (r *AgendamentoPostgresRepository) CriaAgendamento(a *domain.Agendamento) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO agendamentos (
			id,
			cliente_id,
			prestador_id,
			catalogo_id,
			data_hora_inicio,
			data_hora_fim,
			status,
			notas,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`,
		a.ID,
		a.Cliente.ID,
		a.Prestador.ID,
		a.Catalogo.ID,
		a.DataHoraInicio,
		a.DataHoraFim,
		a.Status,
		a.Notas,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AgendamentoPostgresRepository) BuscarPorPrestadorEPeriodo(prestadorID string, inicio time.Time, fim time.Time) ([]*domain.Agendamento, error) {

	query := `
	SELECT
		a.id,
		a.data_hora_inicio,
		a.data_hora_fim,
		a.status,
		a.notas,

		c.id, c.nome, c.email, c.telefone,
		p.id, p.nome, p.cpf, p.email, p.telefone,
		cat.id, cat.nome, cat.duracao_padrao, cat.preco, cat.categoria
	FROM agendamentos a
	JOIN clientes c   ON c.id = a.cliente_id
	JOIN prestadores p ON p.id = a.prestador_id
	JOIN catalogos cat ON cat.id = a.catalogo_id
	WHERE a.prestador_id = $1
	  AND a.data_hora_inicio >= $2
	  AND a.data_hora_fim    <= $3
	ORDER BY a.data_hora_inicio
	`

	rows, err := r.db.Query(query, prestadorID, inicio, fim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendamentos []*domain.Agendamento

	for rows.Next() {
		var a domain.Agendamento
		var cliente domain.Cliente
		var prestador domain.Prestador
		var catalogo domain.Catalogo

		err := rows.Scan(
			&a.ID,
			&a.DataHoraInicio,
			&a.DataHoraFim,
			&a.Status,
			&a.Notas,

			&cliente.ID,
			&cliente.Nome,
			&cliente.Email,
			&cliente.Telefone,

			&prestador.ID,
			&prestador.Nome,
			&prestador.Cpf,
			&prestador.Email,
			&prestador.Telefone,

			&catalogo.ID,
			&catalogo.Nome,
			&catalogo.DuracaoPadrao,
			&catalogo.Preco,
			&catalogo.Categoria,
		)
		if err != nil {
			return nil, err
		}

		a.Cliente = &cliente
		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	return agendamentos, nil
}

func (r *AgendamentoPostgresRepository) BuscarPorClienteEPeriodo(clienteID string, inicio time.Time, fim time.Time) ([]*domain.Agendamento, error) {

	query := `
	SELECT
		a.id,
		a.data_hora_inicio,
		a.data_hora_fim,
		a.status,
		a.notas,

		p.id, p.nome, p.cpf, p.email, p.telefone,
		cat.id, cat.nome, cat.duracao_padrao, cat.preco, cat.categoria
	FROM agendamentos a
	JOIN prestadores p ON p.id = a.prestador_id
	JOIN catalogos cat ON cat.id = a.catalogo_id
	WHERE a.cliente_id = $1
	  AND a.data_hora_inicio >= $2
	  AND a.data_hora_fim    <= $3
	ORDER BY a.data_hora_inicio
	`

	rows, err := r.db.Query(query, clienteID, inicio, fim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agendamentos []*domain.Agendamento

	for rows.Next() {
		var a domain.Agendamento
		var prestador domain.Prestador
		var catalogo domain.Catalogo

		err := rows.Scan(
			&a.ID,
			&a.DataHoraInicio,
			&a.DataHoraFim,
			&a.Status,
			&a.Notas,

			&prestador.ID,
			&prestador.Nome,
			&prestador.Cpf,
			&prestador.Email,
			&prestador.Telefone,

			&catalogo.ID,
			&catalogo.Nome,
			&catalogo.DuracaoPadrao,
			&catalogo.Preco,
			&catalogo.Categoria,
		)
		if err != nil {
			return nil, err
		}

		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	return agendamentos, nil
}
