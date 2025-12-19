package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

func FromCatalogo(c *domain.Catalogo) *output.CatalogoOutput {
	return &output.CatalogoOutput{
		ID:            c.ID,
		Nome:          c.Nome,
		DuracaoPadrao: c.DuracaoPadrao,
		Preco:         c.Preco,
		Categoria:     c.Categoria,
	}
}
