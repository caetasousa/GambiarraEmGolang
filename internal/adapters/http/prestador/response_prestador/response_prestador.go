package response_prestador

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorPostResponse struct {
	ID       string                               `json:"id"`
	Nome     string                               `json:"nome"`
	Email    string                               `json:"email"`
	Telefone string                               `json:"telefone"`
	Ativo    bool                                 `json:"ativo"`
	Catalogo []response_catalogo.CatalogoResponse `json:"catalogo"`
}

func FromPostPrestador(p *domain.Prestador) PrestadorPostResponse {
	catalogo := make([]response_catalogo.CatalogoResponse, len(p.Catalogo))
	for i, c := range p.Catalogo {
		catalogo[i] = response_catalogo.FromCatalogo(&c)
	}

	return PrestadorPostResponse{
		ID:       p.ID,
		Nome:     p.Nome,
		Email:    p.Email,
		Telefone: p.Telefone,
		Ativo:    p.Ativo,
		Catalogo: catalogo,
	}
}
