package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	// Inicia uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // rollback automático se ocorrer erro

	// 1️⃣ Insere prestador
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
		return err
	}

	// 2️⃣ Vincula catálogos
	if len(prestador.Catalogo) > 0 {
		// ON CONFLICT precisa do nome da constraint ou coluna
		stmt, err := tx.Prepare(`
			INSERT INTO prestador_catalogos (prestador_id, catalogo_id)
			VALUES ($1, $2)
			ON CONFLICT (prestador_id, catalogo_id) DO NOTHING
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, catalogo := range prestador.Catalogo {
			if _, err := stmt.Exec(prestador.ID, catalogo.ID); err != nil {
				log.Printf("erro ao inserir prestador_catalogos: prestador_id=%s, catalogo_id=%s, erro=%v",
					prestador.ID, catalogo.ID, err)
				return err
			}
		}
	}

	// 3️⃣ Commit da transação
	if err := tx.Commit(); err != nil {
		log.Printf("erro ao fazer commit: %v", err)
		return err
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
		-- Dados do Catálogo
		c.id AS catalogo_id,
		c.nome AS catalogo_nome,
		c.duracao_padrao AS catalogo_duracao_padrao,
		c.preco AS catalogo_preco,
		c.imagem_url AS catalogo_imagem_url,
		c.categoria AS catalogo_categoria,
		-- Dados da Agenda Diária
		ad.id AS agenda_id,
		ad.data AS agenda_data,
		-- Dados dos Intervalos Diários
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
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var prestador *domain.Prestador
	catalogosMap := make(map[string]*domain.Catalogo)
	agendasMap := make(map[string]*domain.AgendaDiaria)
	intervalosMap := make(map[string]bool) // ✅ NOVO: Controle de intervalos duplicados

	for rows.Next() {
		var (
			// Prestador
			pID, pNome, pCpf, pEmail, pTelefone, pImagemUrl string
			pAtivo                                          bool

			// Catálogo (nullable)
			catalogoID            sql.NullString
			catalogoNome          sql.NullString
			catalogoDuracaoPadrao sql.NullInt64
			catalogoPreco         sql.NullInt64
			catalogoImagemUrl     sql.NullString
			catalogoCategoria     sql.NullString

			// Agenda (nullable)
			agendaID   sql.NullString
			agendaData sql.NullTime

			// Intervalo (nullable)
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
			return nil, fmt.Errorf("erro ao fazer scan: %w", err)
		}

		// Inicializa prestador apenas uma vez
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

		// Adiciona catálogo se existir e ainda não foi adicionado
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

		// Processa agenda e intervalos
		if agendaID.Valid {
			// Adiciona agenda se ainda não existe
			if _, exists := agendasMap[agendaID.String]; !exists {
				agendasMap[agendaID.String] = &domain.AgendaDiaria{
					Id:         agendaID.String,
					Data:       agendaData.Time.Format("2006-01-02"),
					Intervalos: []domain.IntervaloDiario{},
				}
			}

			// ✅ Adiciona intervalo se existir E ainda não foi adicionado
			if intervaloID.Valid {
				// Cria chave única: agendaID + intervaloID
				chaveIntervalo := fmt.Sprintf("%s:%s", agendaID.String, intervaloID.String)
				
				// ✅ Verifica se já foi adicionado
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
					
					// ✅ Marca como adicionado
					intervalosMap[chaveIntervalo] = true
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar rows: %w", err)
	}

	// Se nenhum resultado foi encontrado
	if prestador == nil {
		return nil, sql.ErrNoRows
	}

	// Converte maps para slices
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

	if err == sql.ErrNoRows {
		// CPF não existe → não é erro, retorna nil
		return nil, nil
	}
	if err != nil {
		return nil, err
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
		return nil, err
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
			return nil, err
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
		return nil, nil
	}

	return agenda, nil
}

func (r *PrestadorPostgresRepository) Atualizar(input *input.AlterarPrestadorInput) error {
	// Inicia uma transação
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// 1️⃣ Atualiza dados do prestador (APENAS campos editáveis - sem ID e CPF)
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
		return fmt.Errorf("erro ao atualizar prestador: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // ✅ Retorna erro padrão
	}

	// 2️⃣ Remove todos os catálogos antigos
	_, err = tx.Exec(`
		DELETE FROM prestador_catalogos WHERE prestador_id = $1
	`, input.Id)
	if err != nil {
		return fmt.Errorf("erro ao remover catálogos antigos: %w", err)
	}

	// 3️⃣ Insere novos catálogos
	if len(input.CatalogoIDs) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO prestador_catalogos (prestador_id, catalogo_id)
			VALUES ($1, $2)
		`)
		if err != nil {
			return fmt.Errorf("erro ao preparar inserção de catálogos: %w", err)
		}
		defer stmt.Close()

		for _, catalogoID := range input.CatalogoIDs {
			_, err := stmt.Exec(input.Id, catalogoID)
			if err != nil {
				// Detecta erro de FK (catálogo não existe)
				if pqErr, ok := err.(*pq.Error); ok {
					if pqErr.Code == "23503" { // foreign_key_violation
						return fmt.Errorf("catálogo %s não existe", catalogoID)
					}
				}
				// Alternativa sem lib pq
				if strings.Contains(err.Error(), "foreign key") ||
					strings.Contains(err.Error(), "violates") {
					return fmt.Errorf("catálogo %s não existe", catalogoID)
				}
				return fmt.Errorf("erro ao inserir catálogo %s: %w", catalogoID, err)
			}
		}
	}

	// 4️⃣ Commit da transação
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
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
		-- Dados do Catálogo
		c.id AS catalogo_id,
		c.nome AS catalogo_nome,
		c.duracao_padrao AS catalogo_duracao_padrao,
		c.preco AS catalogo_preco,
		c.imagem_url AS catalogo_imagem_url,
		c.categoria AS catalogo_categoria,
		-- Dados da Agenda Diária
		ad.id AS agenda_id,
		ad.data AS agenda_data,
		-- Dados dos Intervalos Diários
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
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	// Mantém a ordem dos prestadores
	var prestadoresOrdenados []string
	prestadoresMap := make(map[string]*domain.Prestador)
	catalogosMap := make(map[string]map[string]*domain.Catalogo)
	agendasMap := make(map[string]map[string]*domain.AgendaDiaria)
	intervalosMap := make(map[string]map[string]bool) // ✅ NOVO: Controle de intervalos duplicados

	for rows.Next() {
		var (
			// Prestador
			pID, pNome, pCpf, pEmail, pTelefone, pImagemUrl string
			pAtivo                                          bool

			// Catálogo (nullable devido ao LEFT JOIN)
			catalogoID            sql.NullString
			catalogoNome          sql.NullString
			catalogoDuracaoPadrao sql.NullInt64
			catalogoPreco         sql.NullInt64
			catalogoImagemUrl     sql.NullString
			catalogoCategoria     sql.NullString

			// Agenda (nullable devido ao LEFT JOIN)
			agendaID   sql.NullString
			agendaData sql.NullTime

			// Intervalo (nullable devido ao LEFT JOIN)
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
			return nil, fmt.Errorf("erro ao fazer scan: %w", err)
		}

		// Inicializa prestador se não existe
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
			intervalosMap[pID] = make(map[string]bool) // ✅ Inicializa controle de intervalos
		}

		// Adiciona catálogo se existir e ainda não foi adicionado
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

		// Processa agenda e intervalos
		if agendaID.Valid {
			// Adiciona agenda se ainda não existe
			if _, exists := agendasMap[pID][agendaID.String]; !exists {
				agendasMap[pID][agendaID.String] = &domain.AgendaDiaria{
					Id:         agendaID.String,
					Data:       agendaData.Time.Format("2006-01-02"),
					Intervalos: []domain.IntervaloDiario{},
				}
			}

			// ✅ Adiciona intervalo se existir E ainda não foi adicionado
			if intervaloID.Valid {
				// Cria chave única: agendaID + intervaloID
				chaveIntervalo := fmt.Sprintf("%s:%s", agendaID.String, intervaloID.String)
				
				// ✅ Verifica se já foi adicionado
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
					
					// ✅ Marca como adicionado
					intervalosMap[pID][chaveIntervalo] = true
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar rows: %w", err)
	}

	// Converte para slice mantendo a ordem
	prestadores := make([]*domain.Prestador, 0, len(prestadoresOrdenados))
	for _, prestadorID := range prestadoresOrdenados {
		prestador := prestadoresMap[prestadorID]

		// Adiciona catálogos ao prestador
		for _, cat := range catalogosMap[prestadorID] {
			prestador.Catalogo = append(prestador.Catalogo, *cat)
		}

		// Adiciona agendas ao prestador
		for _, agenda := range agendasMap[prestadorID] {
			prestador.Agenda = append(prestador.Agenda, *agenda)
		}

		prestadores = append(prestadores, prestador)
	}

	return prestadores, nil
}

// ✅ Contar com filtro obrigatório
func (r *PrestadorPostgresRepository) Contar(ativo bool) (int, error) {
	var total int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM prestadores 
		WHERE ativo = $1
	`, ativo).Scan(&total)
	
	if err != nil {
		return 0, fmt.Errorf("erro ao contar prestadores: %w", err)
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
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AgendaDiariaPostgresRepository) DeletarAgenda(prestadorID string, data string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// 1. Buscar ID da agenda
	var agendaID string
	err = tx.QueryRow(`
		SELECT id FROM agendas_diarias 
		WHERE prestador_id = $1 AND data = $2
	`, prestadorID, data).Scan(&agendaID)
	
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	}
	if err != nil {
		return fmt.Errorf("erro ao buscar agenda: %w", err)
	}

	// 2. Deletar intervalos
	_, err = tx.Exec(`
		DELETE FROM intervalos_diarios 
		WHERE agenda_id = $1
	`, agendaID)
	if err != nil {
		return fmt.Errorf("erro ao deletar intervalos: %w", err)
	}

	// 3. Deletar agenda
	_, err = tx.Exec(`
		DELETE FROM agendas_diarias 
		WHERE id = $1
	`, agendaID)
	if err != nil {
		return fmt.Errorf("erro ao deletar agenda: %w", err)
	}

	return tx.Commit()
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
		-- Dados do Catálogo
		c.id AS catalogo_id,
		c.nome AS catalogo_nome,
		c.duracao_padrao AS catalogo_duracao_padrao,
		c.preco AS catalogo_preco,
		c.imagem_url AS catalogo_imagem_url,
		c.categoria AS catalogo_categoria,
		-- Dados da Agenda Diária
		ad.id AS agenda_id,
		ad.data AS agenda_data,
		-- Dados dos Intervalos Diários
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
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	// Maps para controlar duplicação e manter a ordem
	var prestadoresOrdem []string // Mantém a ordem de inserção
	prestadoresMap := make(map[string]*domain.Prestador)
	catalogosMap := make(map[string]map[string]*domain.Catalogo)
	agendasMap := make(map[string]*domain.AgendaDiaria)
	intervalosMap := make(map[string]bool)

	for rows.Next() {
		var (
			// Prestador
			pID, pNome, pCpf, pEmail, pTelefone, pImagemUrl string
			pAtivo                                          bool

			// Catálogo (nullable devido ao LEFT JOIN)
			catalogoID            sql.NullString
			catalogoNome          sql.NullString
			catalogoDuracaoPadrao sql.NullInt64
			catalogoPreco         sql.NullInt64
			catalogoImagemUrl     sql.NullString
			catalogoCategoria     sql.NullString

			// Agenda (sempre presente devido ao INNER JOIN com a data)
			agendaID   string
			agendaData sql.NullTime

			// Intervalo (nullable devido ao LEFT JOIN)
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
			return nil, fmt.Errorf("erro ao fazer scan: %w", err)
		}

		// Inicializa prestador se não existe
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

		// Adiciona catálogo se existir e ainda não foi adicionado
		if catalogoID.Valid {
			if _, exists := catalogosMap[pID][catalogoID.String]; !exists {
				catalogo := &domain.Catalogo{
					ID:            catalogoID.String,
					Nome:          catalogoNome.String,
					DuracaoPadrao: int(catalogoDuracaoPadrao.Int64),
					Preco:         int(catalogoPreco.Int64),
					Categoria:     catalogoCategoria.String,
					ImagemUrl:     "", // Inicializa vazio
				}
				if catalogoImagemUrl.Valid {
					catalogo.ImagemUrl = catalogoImagemUrl.String
				}
				catalogosMap[pID][catalogoID.String] = catalogo
			}
		}

		// Adiciona agenda se ainda não existe (será sempre a mesma para este prestador)
		if _, exists := agendasMap[pID]; !exists {
			agendasMap[pID] = &domain.AgendaDiaria{
				Id:         agendaID,
				Data:       agendaData.Time.Format("2006-01-02"),
				Intervalos: []domain.IntervaloDiario{},
			}
		}

		// Adiciona intervalo se existir e ainda não foi adicionado
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
		return nil, fmt.Errorf("erro ao iterar rows: %w", err)
	}

	// Se nenhum prestador foi encontrado
	if len(prestadoresMap) == 0 {
		return []*domain.Prestador{}, nil
	}

	// Converte maps para slice final mantendo a ordem
	prestadores := make([]*domain.Prestador, 0, len(prestadoresOrdem))
	for _, prestadorID := range prestadoresOrdem {
		prestador := prestadoresMap[prestadorID]

		// Adiciona catálogos ao prestador
		for _, cat := range catalogosMap[prestadorID] {
			prestador.Catalogo = append(prestador.Catalogo, *cat)
		}

		// Adiciona agenda ao prestador
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
		return 0, fmt.Errorf("erro ao contar prestadores disponíveis: %w", err)
	}
	
	return total, nil
}