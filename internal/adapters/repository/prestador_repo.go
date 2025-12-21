package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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
		INSERT INTO prestadores (id, nome, cpf, email, telefone, ativo, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`,
		prestador.ID,
		prestador.Nome,
		prestador.Cpf,
		prestador.Email,
		prestador.Telefone,
		prestador.Ativo,
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
		log.Printf("Prestador %s inserido com sucesso", prestador.ID)
		return err
	}

	return nil
}

func (r *PrestadorPostgresRepository) BuscarPorId(id string) (*domain.Prestador, error) {
	var prestador domain.Prestador
	var catalogosJSON, agendasJSON []byte

	query := `SELECT 
        p.id,
        p.nome,
        p.cpf,
        p.email,
        p.telefone,
        p.ativo,
        COALESCE(
            json_agg(DISTINCT jsonb_build_object(
                'id', c.id,
                'nome', c.nome,
                'duracao_padrao', c.duracao_padrao,
                'preco', c.preco,
                'categoria', c.categoria
            )) FILTER (WHERE c.id IS NOT NULL),
            '[]'
        ) AS catalogos,
        COALESCE(
            json_agg(DISTINCT jsonb_build_object(
                'id', a.id,
                'data', a.data,
                'intervalos', (
                    SELECT json_agg(jsonb_build_object(
                        'id', i.id,
                        'hora_inicio', i.hora_inicio,
                        'hora_fim', i.hora_fim
                    )) 
                    FROM intervalos_diarios i
                    WHERE i.agenda_id = a.id
                )
            )) FILTER (WHERE a.id IS NOT NULL),
            '[]'
        ) AS agendas
    FROM prestadores p
    LEFT JOIN prestador_catalogos pc ON pc.prestador_id = p.id
    LEFT JOIN catalogos c ON c.id = pc.catalogo_id
    LEFT JOIN agendas_diarias a ON a.prestador_id = p.id
    WHERE p.id = $1
    GROUP BY p.id;`

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&prestador.ID,
		&prestador.Nome,
		&prestador.Cpf,
		&prestador.Email,
		&prestador.Telefone,
		&prestador.Ativo,
		&catalogosJSON,
		&agendasJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(catalogosJSON, &prestador.Catalogo); err != nil {
		return nil, fmt.Errorf("erro ao deserializar catalogos: %w", err)
	}
	if err := json.Unmarshal(agendasJSON, &prestador.Agenda); err != nil {
		return nil, fmt.Errorf("erro ao deserializar agendas: %w", err)
	}

	// Retorna convertendo para seu domain.Prestador
	result := &domain.Prestador{
		ID:       prestador.ID,
		Nome:     prestador.Nome,
		Cpf:      prestador.Cpf,
		Email:    prestador.Email,
		Telefone: prestador.Telefone,
		Ativo:    prestador.Ativo,
		Catalogo: prestador.Catalogo,
		Agenda:   prestador.Agenda,
	}

	return result, nil
}

func (r *PrestadorPostgresRepository) BuscarPorCPF(cpf string) (*domain.Prestador, error) {
	var p domain.Prestador
	err := r.db.QueryRow(`
        SELECT id, nome, cpf, email, telefone, ativo
        FROM prestadores
        WHERE cpf = $1
    `, cpf).Scan(
		&p.ID,
		&p.Nome,
		&p.Cpf,
		&p.Email,
		&p.Telefone,
		&p.Ativo,
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
	// row := r.db.QueryRow(`
	// 	SELECT id, prestador_id, data
	// 	FROM agenda_diaria
	// 	WHERE prestador_id = $1 AND data = $2
	// `, prestadorID, data)

	// var agenda domain.AgendaDiaria
	// if err := row.Scan(
	// 	&agenda.ID,
	// 	&agenda.PrestadorID,
	// 	&agenda.Data,
	// ); err != nil {
	// 	return nil, err
	// }

	//return &agenda, nil
	return nil, nil
}
