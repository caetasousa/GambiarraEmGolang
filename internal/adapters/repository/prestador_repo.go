package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"

	"github.com/lib/pq"
)

type PrestadorPostgresRepository struct {
	db *sql.DB
}

func NewPrestadorPostgresRepository(db *sql.DB) port.PrestadorRepositorio {
	return &PrestadorPostgresRepository{db: db}
}

func (r *PrestadorPostgresRepository) Salvar(prestador *domain.Prestador) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO prestadores (id, nome, cpf, email, telefone, ativo, imagem_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`,
		prestador.ID,
		prestador.Nome,
		prestador.Cpf,
		prestador.Email,
		prestador.Telefone,
		prestador.Ativo,
		prestador.ImagemUrl,
	)
	if err != nil {
		log.Printf("erro ao inserir prestador: %+v", err)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrCPFDuplicado
		}
		return fmt.Errorf("%w: %v", ErrFalhaAoSalvar, err)
	}

	if len(prestador.Catalogo) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO prestador_catalogos (prestador_id, catalogo_id)
			VALUES ($1, $2)
			ON CONFLICT (prestador_id, catalogo_id) DO NOTHING
		`)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
		}
		defer stmt.Close()

		for _, catalogo := range prestador.Catalogo {
			if _, err := stmt.Exec(prestador.ID, catalogo.ID); err != nil {
				log.Printf("erro ao inserir prestador_catalogos: prestador_id=%s, catalogo_id=%s, erro=%v",
					prestador.ID, catalogo.ID, err)
				return fmt.Errorf("%w: %v", ErrFalhaAoSalvar, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("erro ao fazer commit: %v", err)
		return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
	}
	return nil
}

func (r *PrestadorPostgresRepository) BuscarPorId(id string) (*domain.Prestador, error) {
	query := `
	SELECT 
		p.id,
		p.nome,
		p.cpf,
		p.email,
		p.telefone,
		p.ativo,
		p.imagem_url AS prestador_imagem_url,
		c.id AS catalogo_id,
		c.nome AS catalogo_nome,
		c.duracao_padrao AS catalogo_duracao_padrao,
		c.preco AS catalogo_preco,
		c.imagem_url AS catalogo_imagem_url,
		c.categoria AS catalogo_categoria,
		ad.id AS agenda_id,
		ad.data AS agenda_data,
		id.id AS intervalo_id,
		id.hora_inicio AS intervalo_hora_inicio,
		id.hora_fim AS intervalo_hora_fim
	FROM prestadores p
	LEFT JOIN prestador_catalogos pc ON p.id = pc.prestador_id
	LEFT JOIN catalogos c ON pc.catalogo_id = c.id
	LEFT JOIN agendas_diarias ad ON p.id = ad.prestador_id
	LEFT JOIN intervalos_diarios id ON ad.id = id.agenda_id
	WHERE p.id = $1
	ORDER BY 
		c.nome,
		ad.data,
		id.hora_inicio
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}
	defer rows.Close()

	var prestador *domain.Prestador
	catalogosMap := make(map[string]*domain.Catalogo)
	agendasMap := make(map[string]*domain.AgendaDiaria)
	intervalosMap := make(map[string]bool)

	for rows.Next() {
		var (
			pID, pNome, pCpf, pEmail, pTelefone, pImagemUrl string
			pAtivo                                          bool

			catalogoID            sql.NullString
			catalogoNome          sql.NullString
			catalogoDuracaoPadrao sql.NullInt64
			catalogoPreco         sql.NullInt64
			catalogoImagemUrl     sql.NullString
			catalogoCategoria     sql.NullString

			agendaID   sql.NullString
			agendaData sql.NullTime

			intervaloID         sql.NullString
			intervaloHoraInicio sql.NullTime
			intervaloHoraFim    sql.NullTime
		)

		err := rows.Scan(
			&pID, &pNome, &pCpf, &pEmail, &pTelefone, &pAtivo, &pImagemUrl,
			&catalogoID, &catalogoNome, &catalogoDuracaoPadrao, &catalogoPreco,
			&catalogoImagemUrl, &catalogoCategoria,
			&agendaID, &agendaData,
			&intervaloID, &intervaloHoraInicio, &intervaloHoraFim,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		if prestador == nil {
			prestador = &domain.Prestador{
				ID:        pID,
				Nome:      pNome,
				Cpf:       pCpf,
				Email:     pEmail,
				Telefone:  pTelefone,
				Ativo:     pAtivo,
				ImagemUrl: pImagemUrl,
			}
		}

		if catalogoID.Valid {
			if _, exists := catalogosMap[catalogoID.String]; !exists {
				catalogo := &domain.Catalogo{
					ID:            catalogoID.String,
					Nome:          catalogoNome.String,
					DuracaoPadrao: int(catalogoDuracaoPadrao.Int64),
					Preco:         int(catalogoPreco.Int64),
					Categoria:     catalogoCategoria.String,
				}
				if catalogoImagemUrl.Valid {
					catalogo.ImagemUrl = catalogoImagemUrl.String
				}
				catalogosMap[catalogoID.String] = catalogo
			}
		}

		if agendaID.Valid {
			if _, exists := agendasMap[agendaID.String]; !exists {
				agendasMap[agendaID.String] = &domain.AgendaDiaria{
					Id:         agendaID.String,
					Data:       agendaData.Time.Format("2006-01-02"),
					Intervalos: []domain.IntervaloDiario{},
				}
			}

			if intervaloID.Valid {
				chaveIntervalo := fmt.Sprintf("%s:%s", agendaID.String, intervaloID.String)

				if !intervalosMap[chaveIntervalo] {
					intervalo := domain.IntervaloDiario{
						Id:         intervaloID.String,
						HoraInicio: intervaloHoraInicio.Time,
						HoraFim:    intervaloHoraFim.Time,
					}
					agendasMap[agendaID.String].Intervalos = append(
						agendasMap[agendaID.String].Intervalos,
						intervalo,
					)

					intervalosMap[chaveIntervalo] = true
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	if prestador == nil {
		return nil, ErrPrestadorNaoEncontrado
	}

	prestador.Catalogo = make([]domain.Catalogo, 0, len(catalogosMap))
	for _, cat := range catalogosMap {
		prestador.Catalogo = append(prestador.Catalogo, *cat)
	}

	prestador.Agenda = make([]domain.AgendaDiaria, 0, len(agendasMap))
	for _, agenda := range agendasMap {
		prestador.Agenda = append(prestador.Agenda, *agenda)
	}

	return prestador, nil
}

func (r *PrestadorPostgresRepository) BuscarPorCPF(cpf string) (*domain.Prestador, error) {
	var p domain.Prestador
	err := r.db.QueryRow(`
        SELECT id, nome, cpf, email, telefone, ativo, imagem_url
        FROM prestadores
        WHERE cpf = $1
    `, cpf).Scan(
		&p.ID,
		&p.Nome,
		&p.Cpf,
		&p.Email,
		&p.Telefone,
		&p.Ativo,
		&p.ImagemUrl,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // CPF não existe → permitido para novo cadastro
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return &p, nil
}

func (r *PrestadorPostgresRepository) BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error) {
	rows, err := r.db.Query(`
		SELECT
			a.id,
			a.data,
			i.id,
			(a.data + i.hora_inicio) AT TIME ZONE 'UTC',
			(a.data + i.hora_fim)    AT TIME ZONE 'UTC'
		FROM agendas_diarias a
		LEFT JOIN intervalos_diarios i ON i.agenda_id = a.id
		WHERE a.prestador_id = $1
		  AND a.data = $2
		ORDER BY i.hora_inicio
	`, prestadorID, data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}
	defer rows.Close()

	var agenda *domain.AgendaDiaria

	for rows.Next() {
		var (
			agendaID    string
			dataAgenda  time.Time
			intervaloID sql.NullString
			horaInicio  sql.NullTime
			horaFim     sql.NullTime
		)

		if err := rows.Scan(
			&agendaID,
			&dataAgenda,
			&intervaloID,
			&horaInicio,
			&horaFim,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		if agenda == nil {
			agenda = &domain.AgendaDiaria{
				Id:   agendaID,
				Data: dataAgenda.Format("2006-01-02"),
			}
		}

		if intervaloID.Valid {
			agenda.Intervalos = append(agenda.Intervalos, domain.IntervaloDiario{
				Id:         intervaloID.String,
				HoraInicio: horaInicio.Time,
				HoraFim:    horaFim.Time,
			})
		}
	}

	if agenda == nil {
		return nil, nil // Agenda não existe para esta data
	}

	return agenda, nil
}

func (r *PrestadorPostgresRepository) Atualizar(input *input.AlterarPrestadorInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE prestadores 
		SET 
			nome = $1,
			email = $2,
			telefone = $3,
			imagem_url = $4
		WHERE id = $5
	`,
		input.Nome,
		input.Email,
		input.Telefone,
		input.ImagemUrl,
		input.Id,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
	}

	if rowsAffected == 0 {
		return ErrPrestadorNaoEncontrado
	}

	_, err = tx.Exec(`
		DELETE FROM prestador_catalogos WHERE prestador_id = $1
	`, input.Id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
	}

	if len(input.CatalogoIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO prestador_catalogos (prestador_id, catalogo_id)
			VALUES ($1, $2)
		`)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
		}
		defer stmt.Close()

		for _, catalogoID := range input.CatalogoIDs {
			_, err := stmt.Exec(input.Id, catalogoID)
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
					return fmt.Errorf("%w: %s", ErrCatalogoNaoEncontrado, catalogoID)
				}
				return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaDeTransacao, err)
	}

	return nil
}

func (r *PrestadorPostgresRepository) Listar(input *input.PrestadorListInput) ([]*domain.Prestador, error) {
	offset := (input.Page - 1) * input.Limit

	query := `
	WITH prestadores_paginados AS (
		SELECT 
			id, nome, cpf, email, telefone, ativo, imagem_url, created_at
		FROM prestadores
		WHERE ativo = $3
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	)
	SELECT 
		p.id,
		p.nome,
		p.cpf,
		p.email,
		p.telefone,
		p.ativo,
		p.imagem_url,
		c.id AS catalogo_id,
		c.nome AS catalogo_nome,
		c.duracao_padrao AS catalogo_duracao_padrao,
		c.preco AS catalogo_preco,
		c.imagem_url AS catalogo_imagem_url,
		c.categoria AS catalogo_categoria,
		ad.id AS agenda_id,
		ad.data AS agenda_data,
		id.id AS intervalo_id,
		id.hora_inicio AS intervalo_hora_inicio,
		id.hora_fim AS intervalo_hora_fim
	FROM prestadores_paginados p
	LEFT JOIN prestador_catalogos pc ON p.id = pc.prestador_id
	LEFT JOIN catalogos c ON pc.catalogo_id = c.id
	LEFT JOIN agendas_diarias ad ON p.id = ad.prestador_id
	LEFT JOIN intervalos_diarios id ON ad.id = id.agenda_id
	ORDER BY 
		p.created_at DESC,
		c.nome NULLS LAST,
		ad.data NULLS LAST,
		id.hora_inicio NULLS LAST
	`

	rows, err := r.db.Query(query, input.Limit, offset, input.Ativo)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}
	defer rows.Close()

	var prestadoresOrdenados []string
	prestadoresMap := make(map[string]*domain.Prestador)
	catalogosMap := make(map[string]map[string]*domain.Catalogo)
	agendasMap := make(map[string]map[string]*domain.AgendaDiaria)
	intervalosMap := make(map[string]map[string]bool)

	for rows.Next() {
		var (
			pID, pNome, pCpf, pEmail, pTelefone, pImagemUrl string
			pAtivo                                          bool

			catalogoID            sql.NullString
			catalogoNome          sql.NullString
			catalogoDuracaoPadrao sql.NullInt64
			catalogoPreco         sql.NullInt64
			catalogoImagemUrl     sql.NullString
			catalogoCategoria     sql.NullString

			agendaID   sql.NullString
			agendaData sql.NullTime

			intervaloID         sql.NullString
			intervaloHoraInicio sql.NullTime
			intervaloHoraFim    sql.NullTime
		)

		err := rows.Scan(
			&pID, &pNome, &pCpf, &pEmail, &pTelefone, &pAtivo, &pImagemUrl,
			&catalogoID, &catalogoNome, &catalogoDuracaoPadrao, &catalogoPreco,
			&catalogoImagemUrl, &catalogoCategoria,
			&agendaID, &agendaData,
			&intervaloID, &intervaloHoraInicio, &intervaloHoraFim,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		if _, exists := prestadoresMap[pID]; !exists {
			prestadoresOrdenados = append(prestadoresOrdenados, pID)
			prestadoresMap[pID] = &domain.Prestador{
				ID:        pID,
				Nome:      pNome,
				Cpf:       pCpf,
				Email:     pEmail,
				Telefone:  pTelefone,
				Ativo:     pAtivo,
				ImagemUrl: pImagemUrl,
				Catalogo:  []domain.Catalogo{},
				Agenda:    []domain.AgendaDiaria{},
			}
			catalogosMap[pID] = make(map[string]*domain.Catalogo)
			agendasMap[pID] = make(map[string]*domain.AgendaDiaria)
			intervalosMap[pID] = make(map[string]bool)
		}

		if catalogoID.Valid {
			if _, exists := catalogosMap[pID][catalogoID.String]; !exists {
				catalogo := &domain.Catalogo{
					ID:            catalogoID.String,
					Nome:          catalogoNome.String,
					DuracaoPadrao: int(catalogoDuracaoPadrao.Int64),
					Preco:         int(catalogoPreco.Int64),
					Categoria:     catalogoCategoria.String,
					ImagemUrl:     "",
				}
				if catalogoImagemUrl.Valid {
					catalogo.ImagemUrl = catalogoImagemUrl.String
				}
				catalogosMap[pID][catalogoID.String] = catalogo
			}
		}

		if agendaID.Valid {
			if _, exists := agendasMap[pID][agendaID.String]; !exists {
				agendasMap[pID][agendaID.String] = &domain.AgendaDiaria{
					Id:         agendaID.String,
					Data:       agendaData.Time.Format("2006-01-02"),
					Intervalos: []domain.IntervaloDiario{},
				}
			}

			if intervaloID.Valid {
				chaveIntervalo := fmt.Sprintf("%s:%s", agendaID.String, intervaloID.String)

				if !intervalosMap[pID][chaveIntervalo] {
					intervalo := domain.IntervaloDiario{
						Id:         intervaloID.String,
						HoraInicio: intervaloHoraInicio.Time,
						HoraFim:    intervaloHoraFim.Time,
					}
					agendasMap[pID][agendaID.String].Intervalos = append(
						agendasMap[pID][agendaID.String].Intervalos,
						intervalo,
					)

					intervalosMap[pID][chaveIntervalo] = true
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	prestadores := make([]*domain.Prestador, 0, len(prestadoresOrdenados))
	for _, prestadorID := range prestadoresOrdenados {
		prestador := prestadoresMap[prestadorID]

		for _, cat := range catalogosMap[prestadorID] {
			prestador.Catalogo = append(prestador.Catalogo, *cat)
		}

		for _, agenda := range agendasMap[prestadorID] {
			prestador.Agenda = append(prestador.Agenda, *agenda)
		}

		prestadores = append(prestadores, prestador)
	}

	return prestadores, nil
}

func (r *PrestadorPostgresRepository) Contar(ativo bool) (int, error) {
	var total int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM prestadores 
		WHERE ativo = $1
	`, ativo).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return total, nil
}

func (r *PrestadorPostgresRepository) AtualizarStatus(id string, ativo bool) error {
	result, err := r.db.Exec(`
		UPDATE prestadores 
		SET ativo = $1
		WHERE id = $2
	`, ativo, id)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFalhaAoAtualizar, err)
	}

	if rowsAffected == 0 {
		return ErrPrestadorNaoEncontrado
	}

	return nil
}

func (r *PrestadorPostgresRepository) BuscarPrestadoresDisponiveisPorData(data string, page, limit int) ([]*domain.Prestador, error) {
	offset := (page - 1) * limit

	query := `
	WITH prestadores_paginados AS (
		SELECT DISTINCT
			p.id, p.nome, p.cpf, p.email, p.telefone, p.ativo, p.imagem_url, p.created_at
		FROM prestadores p
		INNER JOIN agendas_diarias ad ON p.id = ad.prestador_id AND ad.data = $1
		WHERE p.ativo = TRUE
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	)
	SELECT 
		p.id,
		p.nome,
		p.cpf,
		p.email,
		p.telefone,
		p.ativo,
		p.imagem_url,
		c.id AS catalogo_id,
		c.nome AS catalogo_nome,
		c.duracao_padrao AS catalogo_duracao_padrao,
		c.preco AS catalogo_preco,
		c.imagem_url AS catalogo_imagem_url,
		c.categoria AS catalogo_categoria,
		ad.id AS agenda_id,
		ad.data AS agenda_data,
		id.id AS intervalo_id,
		id.hora_inicio AS intervalo_hora_inicio,
		id.hora_fim AS intervalo_hora_fim
	FROM prestadores_paginados p
	LEFT JOIN prestador_catalogos pc ON p.id = pc.prestador_id
	LEFT JOIN catalogos c ON pc.catalogo_id = c.id
	LEFT JOIN agendas_diarias ad ON p.id = ad.prestador_id AND ad.data = $1
	LEFT JOIN intervalos_diarios id ON ad.id = id.agenda_id
	ORDER BY 
		p.created_at DESC,
		c.nome NULLS LAST,
		id.hora_inicio NULLS LAST
	`

	rows, err := r.db.Query(query, data, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}
	defer rows.Close()

	var prestadoresOrdem []string
	prestadoresMap := make(map[string]*domain.Prestador)
	catalogosMap := make(map[string]map[string]*domain.Catalogo)
	agendasMap := make(map[string]*domain.AgendaDiaria)
	intervalosMap := make(map[string]bool)

	for rows.Next() {
		var (
			pID, pNome, pCpf, pEmail, pTelefone, pImagemUrl string
			pAtivo                                          bool

			catalogoID            sql.NullString
			catalogoNome          sql.NullString
			catalogoDuracaoPadrao sql.NullInt64
			catalogoPreco         sql.NullInt64
			catalogoImagemUrl     sql.NullString
			catalogoCategoria     sql.NullString

			agendaID   string
			agendaData sql.NullTime

			intervaloID         sql.NullString
			intervaloHoraInicio sql.NullTime
			intervaloHoraFim    sql.NullTime
		)

		err := rows.Scan(
			&pID, &pNome, &pCpf, &pEmail, &pTelefone, &pAtivo, &pImagemUrl,
			&catalogoID, &catalogoNome, &catalogoDuracaoPadrao, &catalogoPreco,
			&catalogoImagemUrl, &catalogoCategoria,
			&agendaID, &agendaData,
			&intervaloID, &intervaloHoraInicio, &intervaloHoraFim,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
		}

		if _, exists := prestadoresMap[pID]; !exists {
			prestadoresOrdem = append(prestadoresOrdem, pID)
			prestadoresMap[pID] = &domain.Prestador{
				ID:        pID,
				Nome:      pNome,
				Cpf:       pCpf,
				Email:     pEmail,
				Telefone:  pTelefone,
				Ativo:     pAtivo,
				ImagemUrl: pImagemUrl,
				Catalogo:  []domain.Catalogo{},
				Agenda:    []domain.AgendaDiaria{},
			}
			catalogosMap[pID] = make(map[string]*domain.Catalogo)
		}

		if catalogoID.Valid {
			if _, exists := catalogosMap[pID][catalogoID.String]; !exists {
				catalogo := &domain.Catalogo{
					ID:            catalogoID.String,
					Nome:          catalogoNome.String,
					DuracaoPadrao: int(catalogoDuracaoPadrao.Int64),
					Preco:         int(catalogoPreco.Int64),
					Categoria:     catalogoCategoria.String,
					ImagemUrl:     "",
				}
				if catalogoImagemUrl.Valid {
					catalogo.ImagemUrl = catalogoImagemUrl.String
				}
				catalogosMap[pID][catalogoID.String] = catalogo
			}
		}

		if _, exists := agendasMap[pID]; !exists {
			agendasMap[pID] = &domain.AgendaDiaria{
				Id:         agendaID,
				Data:       agendaData.Time.Format("2006-01-02"),
				Intervalos: []domain.IntervaloDiario{},
			}
		}

		if intervaloID.Valid {
			chaveIntervalo := fmt.Sprintf("%s:%s", agendaID, intervaloID.String)

			if !intervalosMap[chaveIntervalo] {
				intervalo := domain.IntervaloDiario{
					Id:         intervaloID.String,
					HoraInicio: intervaloHoraInicio.Time,
					HoraFim:    intervaloHoraFim.Time,
				}
				agendasMap[pID].Intervalos = append(
					agendasMap[pID].Intervalos,
					intervalo,
				)
				intervalosMap[chaveIntervalo] = true
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	if len(prestadoresMap) == 0 {
		return []*domain.Prestador{}, nil
	}

	prestadores := make([]*domain.Prestador, 0, len(prestadoresOrdem))
	for _, prestadorID := range prestadoresOrdem {
		prestador := prestadoresMap[prestadorID]

		for _, cat := range catalogosMap[prestadorID] {
			prestador.Catalogo = append(prestador.Catalogo, *cat)
		}

		if agenda, exists := agendasMap[prestadorID]; exists {
			prestador.Agenda = append(prestador.Agenda, *agenda)
		}

		prestadores = append(prestadores, prestador)
	}

	return prestadores, nil
}

func (r *PrestadorPostgresRepository) ContarPrestadoresDisponiveisPorData(data string) (int, error) {
	var total int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT p.id)
		FROM prestadores p
		INNER JOIN agendas_diarias ad ON p.id = ad.prestador_id AND ad.data = $1
		WHERE p.ativo = TRUE
	`, data).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFalhaAoListar, err)
	}

	return total, nil
}
