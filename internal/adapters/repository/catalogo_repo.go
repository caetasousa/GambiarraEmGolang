package repository

import (
	"database/sql"
	"errors"

	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoPostgresRepositorio struct {
	db *sql.DB
}

func NovoCatalogoPostgresRepositorio(db *sql.DB) port.CatalogoRepositorio {
	return &CatalogoPostgresRepositorio{
		db: db,
	}
}

func (r *CatalogoPostgresRepositorio) Salvar(c *domain.Catalogo) error {
	query := `
		INSERT INTO catalogos (id, nome, duracao_padrao, preco, categoria)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		c.ID,
		c.Nome,
		c.DuracaoPadrao,
		c.Preco,
		c.Categoria,
	)

	if err != nil {
		return err
	}

	return nil
}
func (r *CatalogoPostgresRepositorio) BuscarPorId(id string) (*domain.Catalogo, error) {
	query := `
		SELECT id, nome, duracao_padrao, preco, categoria
		FROM catalogos
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var c domain.Catalogo
	err := row.Scan(
		&c.ID,
		&c.Nome,
		&c.DuracaoPadrao,
		&c.Preco,
		&c.Categoria,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
