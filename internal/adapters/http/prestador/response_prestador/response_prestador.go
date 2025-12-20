package response_prestador

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/core/application/output"
)

type PrestadorPostResponse struct {
	ID       string                               `json:"id"`
	Nome     string                               `json:"nome"`
	Email    string                               `json:"email"`
	Telefone string                               `json:"telefone"`
	Ativo    bool                                 `json:"ativo"`
	Catalogo []response_catalogo.CatalogoResponse `json:"catalogo"`
}

func FromCriarPrestadorOutput(o output.CriarPrestadorOutput) PrestadorPostResponse {

	catalogo := make([]response_catalogo.CatalogoResponse, len(o.Catalogo))
	for i, c := range o.Catalogo {
		catalogo[i] = response_catalogo.FromCatalogoResponse(c)
	}

	return PrestadorPostResponse{
		ID:       o.ID,
		Nome:     o.Nome,
		Email:    o.Email,
		Telefone: o.Telefone,
		Ativo:    o.Ativo,
		Catalogo: catalogo,
	}
}
