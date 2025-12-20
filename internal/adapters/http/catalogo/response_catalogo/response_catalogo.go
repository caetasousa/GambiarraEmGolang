package response_catalogo

import "meu-servico-agenda/internal/core/application/output"

type CatalogoResponse struct {
	ID            string `json:"id"`
	Nome          string `json:"nome"`
	DuracaoPadrao int    `json:"duracao_padrao"`
	Preco         int  `json:"preco"`
	Categoria     string `json:"categoria"`
}

func FromCatalogoResponse(o output.CatalogoOutput) CatalogoResponse {
	return CatalogoResponse{
		ID:            o.ID,
		Nome:          o.Nome,
		DuracaoPadrao: o.DuracaoPadrao,
		Preco:         o.Preco,
		Categoria:     o.Categoria,
	}
}
