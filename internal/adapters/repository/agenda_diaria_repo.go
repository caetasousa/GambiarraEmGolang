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
