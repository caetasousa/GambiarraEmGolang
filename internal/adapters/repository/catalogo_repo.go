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
		INSERT INTO catalogos (
			id,
			nome,
			duracao_padrao,
			preco,
			categoria,
			imagem_url
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		query,
		c.ID,
		c.Nome,
		c.DuracaoPadrao,
		c.Preco,
		c.Categoria,
		c.ImagemUrl,
	)

	if err != nil {
		return err
	}

	return nil
}
func (r *CatalogoPostgresRepositorio) BuscarPorId(id string) (*domain.Catalogo, error) {
	query := `
		SELECT id, nome, duracao_padrao, preco, categoria, imagem_url
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
		&c.ImagemUrl,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &c, nil
}

func (r *CatalogoPostgresRepositorio) Listar(limit, offset int) ([]*domain.Catalogo, error) {
	query := `
		SELECT
			id,
			nome,
			duracao_padrao,
			preco,
			categoria,
			imagem_url
		FROM catalogos
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var catalogos []*domain.Catalogo

	for rows.Next() {
		var c domain.Catalogo
		if err := rows.Scan(
			&c.ID,
			&c.Nome,
			&c.DuracaoPadrao,
			&c.Preco,
			&c.Categoria,
			&c.ImagemUrl,
		); err != nil {
			return nil, err
		}
		catalogos = append(catalogos, &c)
	}

	return catalogos, nil
}

func (r *CatalogoPostgresRepositorio) Contar() (int, error) {
	query := `SELECT COUNT(*) FROM catalogos`

	var total int
	err := r.db.QueryRow(query).Scan(&total)
	return total, err
}
