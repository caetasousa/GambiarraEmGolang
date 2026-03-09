package repository

import (
	"database/sql"
	"fmt"
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
		return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
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
		return fmt.Errorf("%w: %v", ErrFalhaAoSalvar, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
	}

	return nil
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
	  AND a.data_hora_inicio < $3
	  AND a.data_hora_fim    > $2
	  AND a.status != 3
	ORDER BY a.data_hora_inicio
	`

	rows, err := r.db.Query(query, prestadorID, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		a.Cliente = &cliente
		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
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

		c.id, c.nome, c.email, c.telefone,
		p.id, p.nome, p.cpf, p.email, p.telefone,
		cat.id, cat.nome, cat.duracao_padrao, cat.preco, cat.categoria
	FROM agendamentos a
	JOIN clientes c   ON c.id = a.cliente_id
	JOIN prestadores p ON p.id = a.prestador_id
	JOIN catalogos cat ON cat.id = a.catalogo_id
	WHERE a.cliente_id = $1
	  AND a.data_hora_inicio < $3
	  AND a.data_hora_fim    > $2
	  AND a.status != 3
	ORDER BY a.data_hora_inicio
	`

	rows, err := r.db.Query(query, clienteID, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		a.Cliente = &cliente
		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return agendamentos, nil
}

func (r *AgendamentoPostgresRepository) BuscarAgendamentoClienteAPartirDaData(clienteID string, data time.Time) ([]*domain.Agendamento, error) {
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
	WHERE a.cliente_id = $1
	  AND a.data_hora_inicio >= $2
	ORDER BY a.data_hora_inicio
	`

	rows, err := r.db.Query(query, clienteID, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		a.Cliente = &cliente
		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return agendamentos, nil
}

func (r *AgendamentoPostgresRepository) BuscarAgendamentoPrestadorAPartirDaData(prestadorID string, data time.Time) ([]*domain.Agendamento, error) {
	query := `
	SELECT
		a.id,
		a.data_hora_inicio,
		a.data_hora_fim,
		a.status,
		a.notas,

		c.id, c.nome, c.email, c.telefone,
		p.id, p.nome, p.cpf, p.email, p.telefone, p.ativo, p.imagem_url,
		cat.id, cat.nome, cat.duracao_padrao, cat.preco, cat.imagem_url, cat.categoria

	FROM agendamentos a
	JOIN clientes c ON c.id = a.cliente_id
	JOIN prestadores p ON p.id = a.prestador_id
	JOIN catalogos cat ON cat.id = a.catalogo_id
	WHERE a.prestador_id = $1
	  AND a.data_hora_inicio >= $2
	ORDER BY a.data_hora_inicio
	`

	rows, err := r.db.Query(query, prestadorID, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}
	defer rows.Close()

	var agendamentos []*domain.Agendamento

	for rows.Next() {
		var agendamentoID, clienteID, clienteNome, clienteEmail, clienteTelefone string
		var pID, prestadorNome, prestadorEmail, prestadorTelefone string
		var prestadorCpf domain.CPF
		var prestadorAtivo bool
		var prestadorImagemUrl, catalogoImagemUrl sql.NullString
		var catalogoID, catalogoNome, catalogoCategoria string
		var catalogoDuracao, catalogoPreco int
		var dataHoraInicio, dataHoraFim time.Time
		var status int
		var notas sql.NullString

		err := rows.Scan(
			&agendamentoID,
			&dataHoraInicio,
			&dataHoraFim,
			&status,
			&notas,

			&clienteID,
			&clienteNome,
			&clienteEmail,
			&clienteTelefone,

			&pID,
			&prestadorNome,
			&prestadorCpf,
			&prestadorEmail,
			&prestadorTelefone,
			&prestadorAtivo,
			&prestadorImagemUrl,

			&catalogoID,
			&catalogoNome,
			&catalogoDuracao,
			&catalogoPreco,
			&catalogoImagemUrl,
			&catalogoCategoria,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		cliente := &domain.Cliente{
			ID:       clienteID,
			Nome:     clienteNome,
			Email:    clienteEmail,
			Telefone: clienteTelefone,
		}

		prestador := &domain.Prestador{
			ID:        pID,
			Nome:      prestadorNome,
			Cpf:       prestadorCpf,
			Email:     prestadorEmail,
			Telefone:  prestadorTelefone,
			Ativo:     prestadorAtivo,
			ImagemUrl: prestadorImagemUrl.String,
			Agenda:    []domain.AgendaDiaria{},
		}

		catalogo := &domain.Catalogo{
			ID:            catalogoID,
			Nome:          catalogoNome,
			DuracaoPadrao: catalogoDuracao,
			Preco:         catalogoPreco,
			ImagemUrl:     catalogoImagemUrl.String,
			Categoria:     catalogoCategoria,
		}

		notasStr := ""
		if notas.Valid {
			notasStr = notas.String
		}

		agendamento := &domain.Agendamento{
			ID:             agendamentoID,
			Cliente:        cliente,
			Prestador:      prestador,
			Catalogo:       catalogo,
			DataHoraInicio: dataHoraInicio,
			DataHoraFim:    dataHoraFim,
			Status:         domain.StatusDoAgendamento(status),
			Notas:          notasStr,
		}

		agendamentos = append(agendamentos, agendamento)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return agendamentos, nil
}

func (r *AgendamentoPostgresRepository) ListarPorPrestadorEPeriodoPaginado(prestadorID string, inicio time.Time, fim time.Time, limit, offset int) ([]*domain.Agendamento, error) {
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
	  AND a.data_hora_inicio < $3
	  AND a.data_hora_fim    > $2
	ORDER BY a.data_hora_inicio
	LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(query, prestadorID, inicio, fim, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		a.Cliente = &cliente
		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return agendamentos, nil
}

func (r *AgendamentoPostgresRepository) ListarPorClienteEPeriodoPaginado(clienteID string, inicio time.Time, fim time.Time, limit, offset int) ([]*domain.Agendamento, error) {
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
	WHERE a.cliente_id = $1
	  AND a.data_hora_inicio < $3
	  AND a.data_hora_fim    > $2
	ORDER BY a.data_hora_inicio
	LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(query, clienteID, inicio, fim, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		a.Cliente = &cliente
		a.Prestador = &prestador
		a.Catalogo = &catalogo

		agendamentos = append(agendamentos, &a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return agendamentos, nil
}

func (r *AgendamentoPostgresRepository) ContarPorPrestadorEPeriodo(prestadorID string, inicio time.Time, fim time.Time) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM agendamentos a
	WHERE a.prestador_id = $1
	  AND a.data_hora_inicio < $3
	  AND a.data_hora_fim    > $2
	`

	var total int
	err := r.db.QueryRow(query, prestadorID, inicio, fim).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFalhaAoContar, err)
	}
	return total, nil
}

func (r *AgendamentoPostgresRepository) ContarPorClienteEPeriodo(clienteID string, inicio time.Time, fim time.Time) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM agendamentos a
	WHERE a.cliente_id = $1
	  AND a.data_hora_inicio < $3
	  AND a.data_hora_fim    > $2
	`

	var total int
	err := r.db.QueryRow(query, clienteID, inicio, fim).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFalhaAoContar, err)
	}
	return total, nil
}

func (r *AgendamentoPostgresRepository) BuscarPorID(id string) (*domain.Agendamento, error) {
	query := `
	SELECT
		a.id,
		a.data_hora_inicio,
		a.data_hora_fim,
		a.status,
		a.notas,

		c.id, c.nome, c.email, c.telefone,
		p.id, p.nome, p.cpf, p.email, p.telefone, p.ativo, p.imagem_url,
		cat.id, cat.nome, cat.duracao_padrao, cat.preco, cat.imagem_url, cat.categoria
	FROM agendamentos a
	JOIN clientes c   ON c.id = a.cliente_id
	JOIN prestadores p ON p.id = a.prestador_id
	JOIN catalogos cat ON cat.id = a.catalogo_id
	WHERE a.id = $1
	`

	var agendamentoID, clienteID, clienteNome, clienteEmail, clienteTelefone string
	var pID, prestadorNome, prestadorEmail, prestadorTelefone string
	var prestadorCpf domain.CPF
	var prestadorAtivo bool
	var prestadorImagemUrl, catalogoImagemUrl sql.NullString
	var catalogoID, catalogoNome, catalogoCategoria string
	var catalogoDuracao, catalogoPreco int
	var dataHoraInicio, dataHoraFim time.Time
	var status int
	var notas sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&agendamentoID,
		&dataHoraInicio,
		&dataHoraFim,
		&status,
		&notas,

		&clienteID,
		&clienteNome,
		&clienteEmail,
		&clienteTelefone,

		&pID,
		&prestadorNome,
		&prestadorCpf,
		&prestadorEmail,
		&prestadorTelefone,
		&prestadorAtivo,
		&prestadorImagemUrl,

		&catalogoID,
		&catalogoNome,
		&catalogoDuracao,
		&catalogoPreco,
		&catalogoImagemUrl,
		&catalogoCategoria,
	)
	if err == sql.ErrNoRows {
		return nil, ErrAgendamentoNaoEncontrado
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNaoEncontrado, err)
	}

	notasStr := ""
	if notas.Valid {
		notasStr = notas.String
	}

	return &domain.Agendamento{
		ID: agendamentoID,
		Cliente: &domain.Cliente{
			ID:       clienteID,
			Nome:     clienteNome,
			Email:    clienteEmail,
			Telefone: clienteTelefone,
		},
		Prestador: &domain.Prestador{
			ID:        pID,
			Nome:      prestadorNome,
			Cpf:       prestadorCpf,
			Email:     prestadorEmail,
			Telefone:  prestadorTelefone,
			Ativo:     prestadorAtivo,
			ImagemUrl: prestadorImagemUrl.String,
			Agenda:    []domain.AgendaDiaria{},
		},
		Catalogo: &domain.Catalogo{
			ID:            catalogoID,
			Nome:          catalogoNome,
			DuracaoPadrao: catalogoDuracao,
			Preco:         catalogoPreco,
			ImagemUrl:     catalogoImagemUrl.String,
			Categoria:     catalogoCategoria,
		},
		DataHoraInicio: dataHoraInicio,
		DataHoraFim:    dataHoraFim,
		Status:         domain.StatusDoAgendamento(status),
		Notas:          notasStr,
	}, nil
}

func (r *AgendamentoPostgresRepository) AtualizarStatus(id string, status domain.StatusDoAgendamento) error {
	_, err := r.db.Exec(`UPDATE agendamentos SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
	}
	return nil
}
