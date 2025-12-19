package response_catalogo

import "meu-servico-agenda/internal/core/application/output"

type CatalogoResponse struct {
	ID            string `json:"id"`
	Nome          string `json:"nome"`
	DuracaoPadrao int    `json:"duracao_padrao"`
	Preco         int64  `json:"preco"`
	Categoria     string `json:"categoria"`
}

func FromCatalogoOutput(o output.CatalogoOutput) CatalogoResponse {
	return CatalogoResponse{
		ID:            o.ID,
		Nome:          o.Nome,
		DuracaoPadrao: o.DuracaoPadrao,
		Preco:         int64(o.Preco),
		Categoria:     o.Categoria,
	}
}
