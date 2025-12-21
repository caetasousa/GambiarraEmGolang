package repository

import (
	"database/sql"
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
    // 1️⃣ Busca prestador e seus catálogos
    rows, err := r.db.Query(`
        SELECT 
            p.id, p.nome, p.cpf, p.email, p.telefone, p.ativo,
            c.id, c.nome, c.duracao_padrao, c.preco, c.categoria
        FROM prestadores p
        LEFT JOIN prestador_catalogos pc ON pc.prestador_id = p.id
        LEFT JOIN catalogos c ON c.id = pc.catalogo_id
        WHERE p.id = $1
    `, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var prestador *domain.Prestador
    catalogos := []domain.Catalogo{}

    for rows.Next() {
        var p domain.Prestador
        var c domain.Catalogo
        var cID sql.NullString

        err := rows.Scan(
            &p.ID,
            &p.Nome,
            &p.Cpf,
            &p.Email,
            &p.Telefone,
            &p.Ativo,
            &cID,
            &c.Nome,
            &c.DuracaoPadrao,
            &c.Preco,
            &c.Categoria,
        )
        if err != nil {
            return nil, err
        }

        if prestador == nil {
            prestador = &p
        }

        if cID.Valid {
            c.ID = cID.String
            catalogos = append(catalogos, c)
        }
    }

    if prestador == nil {
        return nil, sql.ErrNoRows
    }

    prestador.Catalogo = catalogos

    // 2️⃣ Busca agendas e intervalos
    agendasRows, err := r.db.Query(`
        SELECT a.id, a.data, i.id, i.hora_inicio, i.hora_fim
        FROM agendas_diarias a
        LEFT JOIN intervalos_diarios i ON i.agenda_id = a.id
        WHERE a.prestador_id = $1
        ORDER BY a.data, i.hora_inicio
    `, id)
    if err != nil {
        return nil, err
    }
    defer agendasRows.Close()

    agendasMap := map[string]*domain.AgendaDiaria{}

    for agendasRows.Next() {
        var aID, iID sql.NullString
        var data string
        var horaInicio, horaFim sql.NullTime

        if err := agendasRows.Scan(&aID, &data, &iID, &horaInicio, &horaFim); err != nil {
            return nil, err
        }

        if !aID.Valid {
            continue
        }

        agenda, ok := agendasMap[aID.String]
        if !ok {
            agenda = &domain.AgendaDiaria{
                Id:         aID.String,
                Data:       data,
                Intervalos: []domain.IntervaloDiario{},
            }
            agendasMap[aID.String] = agenda
        }

        if iID.Valid && horaInicio.Valid && horaFim.Valid {
            agenda.Intervalos = append(agenda.Intervalos, domain.IntervaloDiario{
                Id:         iID.String,
                HoraInicio: horaInicio.Time,
                HoraFim:    horaFim.Time,
            })
        }
    }

    // Converte map para slice
    prestador.Agenda = make([]domain.AgendaDiaria, 0, len(agendasMap))
    for _, agenda := range agendasMap {
        prestador.Agenda = append(prestador.Agenda, *agenda)
    }

    return prestador, nil
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
