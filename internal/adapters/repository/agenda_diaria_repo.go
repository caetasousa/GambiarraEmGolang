package repository

import (
	"database/sql"
	"fmt"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type AgendaDiariaPostgresRepository struct {
	db *sql.DB
}

func NovoAgendaDiariaPostgresRepository(db *sql.DB) port.AgendaDiariaRepositorio {
	return &AgendaDiariaPostgresRepository{db: db}
}

func (r *AgendaDiariaPostgresRepository) Salvar(agenda *domain.AgendaDiaria, prestadorId string) error {
	// Inicia transação
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // rollback automático em caso de erro

	// 1️⃣ Insere agenda
	_, err = tx.Exec(`
		INSERT INTO agendas_diarias (id, prestador_id, data, created_at)
		VALUES ($1, $2, $3, NOW())
	`,
		agenda.Id,
		prestadorId, // precisa ter PrestadorID no domínio
		agenda.Data,
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir agenda: %w", err)
	}

	// 2️⃣ Insere intervalos
	if len(agenda.Intervalos) > 0 {
		stmt, err := tx.Prepare(`
			INSERT INTO intervalos_diarios (id, agenda_id, hora_inicio, hora_fim)
			VALUES ($1, $2, $3, $4)
		`)
		if err != nil {
			return fmt.Errorf("erro ao preparar insert de intervalos: %w", err)
		}
		defer stmt.Close()

		for _, it := range agenda.Intervalos {
			if _, err := stmt.Exec(it.Id, agenda.Id, it.HoraInicio.Format("15:04:05"), it.HoraFim.Format("15:04:05")); err != nil {
				return fmt.Errorf("erro ao inserir intervalo: %w", err)
			}
		}
	}

	// 3️⃣ Commit
	return tx.Commit()
}

func (r *AgendaDiariaPostgresRepository) AtualizarAgenda(agenda *domain.AgendaDiaria, prestadorID string) error {
	// Inicia transação
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	// 0. ✅ VALIDAÇÃO: Verificar se a agenda pertence ao prestador
	var count int
	err = tx.QueryRow(`
		SELECT COUNT(*) 
		FROM agendas_diarias 
		WHERE id = $1 AND prestador_id = $2
	`, agenda.Id, prestadorID).Scan(&count)
	if err != nil {
		return fmt.Errorf("erro ao verificar agenda: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows // Agenda não pertence a este prestador
	}

	// 1. Deletar todos os intervalos antigos
	_, err = tx.Exec(`
		DELETE FROM intervalos_diarios 
		WHERE agenda_id = $1
	`, agenda.Id)
	if err != nil {
		return fmt.Errorf("erro ao deletar intervalos antigos: %w", err)
	}

	// 2. Inserir novos intervalos
	for _, intervalo := range agenda.Intervalos {
		_, err = tx.Exec(`
			INSERT INTO intervalos_diarios (id, agenda_id, hora_inicio, hora_fim)
			VALUES ($1, $2, $3, $4)
		`, intervalo.Id, agenda.Id, intervalo.HoraInicio, intervalo.HoraFim)
		if err != nil {
			return fmt.Errorf("erro ao inserir novo intervalo: %w", err)
		}
	}

	// 3. Commit da transação
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao commitar transação: %w", err)
	}

	return nil
}

func (r *AgendaDiariaPostgresRepository) BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error) {
	query := `
		SELECT 
			ad.id AS agenda_id,
			ad.data AS agenda_data,
			id.id AS intervalo_id,
			id.hora_inicio,
			id.hora_fim
		FROM agendas_diarias ad
		LEFT JOIN intervalos_diarios id ON ad.id = id.agenda_id
		WHERE ad.prestador_id = $1 AND ad.data = $2
		ORDER BY id.hora_inicio
	`

	rows, err := r.db.Query(query, prestadorID, data)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar agenda: %w", err)
	}
	defer rows.Close()

	var agenda *domain.AgendaDiaria
	intervalosMap := make(map[string]bool) // ✅ Controle de duplicação

	for rows.Next() {
		var (
			agendaID   string
			agendaData string

			intervaloID         sql.NullString
			intervaloHoraInicio sql.NullTime
			intervaloHoraFim    sql.NullTime
		)

		err := rows.Scan(
			&agendaID,
			&agendaData,
			&intervaloID,
			&intervaloHoraInicio,
			&intervaloHoraFim,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao fazer scan: %w", err)
		}

		// Inicializa agenda apenas uma vez
		if agenda == nil {
			agenda = &domain.AgendaDiaria{
				Id:         agendaID,
				Data:       agendaData,
				Intervalos: []domain.IntervaloDiario{},
			}
		}

		// Adiciona intervalo se existir e não foi adicionado ainda
		if intervaloID.Valid {
			// ✅ Verifica duplicação
			if !intervalosMap[intervaloID.String] {
				intervalo := domain.IntervaloDiario{
					Id:         intervaloID.String,
					HoraInicio: intervaloHoraInicio.Time,
					HoraFim:    intervaloHoraFim.Time,
				}
				agenda.Intervalos = append(agenda.Intervalos, intervalo)
				intervalosMap[intervaloID.String] = true
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar rows: %w", err)
	}

	// Se não encontrou nenhuma agenda
	if agenda == nil {
		return nil, sql.ErrNoRows
	}

	return agenda, nil
}