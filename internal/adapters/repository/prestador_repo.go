package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
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

	log.Printf("✅ prestador de cpf %s", prestador.Cpf)
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

	log.Printf("Prestador %s inserido com sucesso", prestador.ID)
	return nil
}

func (r *PrestadorPostgresRepository) BuscarPorId(id string) (*domain.Prestador, error) {
	var prestador domain.Prestador
	var catalogosJSON, agendasJSON []byte

	query := `
	SELECT
		p.id,
		p.nome,
		p.cpf,
		p.email,
		p.telefone,
		p.ativo,
		p.imagem_url,
		COALESCE(catalogos.catalogos, '[]') AS catalogos,
		COALESCE(agendas.agendas, '[]') AS agendas
	FROM prestadores p

	-- Catálogos
	LEFT JOIN (
		SELECT
			pc.prestador_id,
			json_agg(
				jsonb_build_object(
					'id', c.id,
					'nome', c.nome,
					'duracao_padrao', c.duracao_padrao,
					'preco', c.preco,
					'categoria', c.categoria
				)
			) AS catalogos
		FROM prestador_catalogos pc
		JOIN catalogos c ON c.id = pc.catalogo_id
		GROUP BY pc.prestador_id
	) catalogos ON catalogos.prestador_id = p.id

	-- Agendas com intervalos
	LEFT JOIN (
		SELECT
			a.prestador_id,
			json_agg(
				jsonb_build_object(
					'id', a.id,
					'data', a.data,
					'intervalos', COALESCE(i.intervalos, '[]')
				)
				ORDER BY a.data
			) AS agendas
		FROM agendas_diarias a
		LEFT JOIN (
			SELECT
				agenda_id,
				json_agg(
					jsonb_build_object(
						'id', id,
						'hora_inicio', hora_inicio,
						'hora_fim', hora_fim
					)
					ORDER BY hora_inicio
				) AS intervalos
			FROM intervalos_diarios
			GROUP BY agenda_id
		) i ON i.agenda_id = a.id
		GROUP BY a.prestador_id
	) agendas ON agendas.prestador_id = p.id

	WHERE p.id = $1;
	`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&prestador.ID,
		&prestador.Nome,
		&prestador.Cpf,
		&prestador.Email,
		&prestador.Telefone,
		&prestador.Ativo,
		&prestador.ImagemUrl,
		&catalogosJSON,
		&agendasJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if err := json.Unmarshal(catalogosJSON, &prestador.Catalogo); err != nil {
		return nil, fmt.Errorf("erro ao deserializar catalogos: %w", err)
	}

	if err := json.Unmarshal(agendasJSON, &prestador.Agenda); err != nil {
		return nil, fmt.Errorf("erro ao deserializar agendas: %w", err)
	}

	return &prestador, nil
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