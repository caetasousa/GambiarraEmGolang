package response_catalogo

import "meu-servico-agenda/internal/core/domain"

type CatalogoResponse struct {
	ID            string `json:"id"`
	Nome          string `json:"nome"`
	DuracaoPadrao int    `json:"duracao_padrao"`
	Preco         int64  `json:"preco"`
	Categoria     string `json:"categoria"`
}

func FromCatalogo(c *domain.Catalogo) CatalogoResponse {
	return CatalogoResponse{
		ID:            c.ID,
		Nome:          c.Nome,
		DuracaoPadrao: c.DuracaoPadrao,
		Preco:         int64(c.Preco),
		Categoria:     c.Categoria,
	}
}
