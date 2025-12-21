package repository

import (
	"database/sql"
	"errors"

	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"

	"github.com/lib/pq"
)

type ClientePostgresRepositorio struct {
	db *sql.DB
}

func NovoClientePostgresRepositorio(db *sql.DB) port.ClienteRepositorio {
	return &ClientePostgresRepositorio{
		db: db,
	}
}

func (r *ClientePostgresRepositorio) Salvar(cliente *domain.Cliente) error {
	query := `
		INSERT INTO clientes (id, nome, email, telefone)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		query,
		cliente.ID,
		cliente.Nome,
		cliente.Email,
		cliente.Telefone,
	)

	if err != nil {
		// trata erro de email duplicado
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("email ja cadastrado")
		}
		return err
	}

	return nil
}

func (r *ClientePostgresRepositorio) BuscarPorId(id string) (*domain.Cliente, error) {
	query := `
		SELECT id, nome, email, telefone
		FROM clientes
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var cliente domain.Cliente
	err := row.Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.Email,
		&cliente.Telefone,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cliente, nil
}
