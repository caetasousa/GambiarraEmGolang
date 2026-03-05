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
		INSERT INTO clientes (id, nome, email, telefone, senha_hash, role)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		query,
		cliente.ID,
		cliente.Nome,
		cliente.Email,
		cliente.Telefone,
		cliente.PasswordHash,
		string(cliente.Role),
	)

	if err != nil {
		// trata erro de email duplicado
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrEmailDuplicado
		}
		return err
	}

	return nil
}

func (r *ClientePostgresRepositorio) BuscarPorId(id string) (*domain.Cliente, error) {
	query := `
		SELECT id, nome, email, telefone, senha_hash, role
		FROM clientes
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var cliente domain.Cliente
	var role string
	err := row.Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.Email,
		&cliente.Telefone,
		&cliente.PasswordHash,
		&role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Cliente não encontrado
		}
		return nil, err
	}

	cliente.Role = domain.Role(role)
	return &cliente, nil
}

func (r *ClientePostgresRepositorio) BuscarPorEmail(email string) (*domain.Cliente, error) {
	query := `
		SELECT id, nome, email, telefone, senha_hash, role
		FROM clientes
		WHERE email = $1
	`

	row := r.db.QueryRow(query, email)

	var cliente domain.Cliente
	var role string
	err := row.Scan(
		&cliente.ID,
		&cliente.Nome,
		&cliente.Email,
		&cliente.Telefone,
		&cliente.PasswordHash,
		&role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Cliente não encontrado
		}
		return nil, err
	}

	cliente.Role = domain.Role(role)
	return &cliente, nil
}
