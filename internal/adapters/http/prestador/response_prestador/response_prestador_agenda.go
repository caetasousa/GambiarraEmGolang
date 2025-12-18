package response_prestador

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorResponse struct {
	ID       string                               `json:"id"`
	Nome     string                               `json:"nome"`
	Email    string                               `json:"email"`
	Telefone string                               `json:"telefone"`
	Ativo    bool                                 `json:"ativo"`
	Catalogo []response_catalogo.CatalogoResponse `json:"catalogo"`
	Agenda   []AgendaDiariaResponse               `json:"agenda"`
}

func FromPrestador(p *domain.Prestador) PrestadorResponse {

	agenda := make([]AgendaDiariaResponse, len(p.Agenda))
	for i, ag := range p.Agenda {
		agenda[i] = FromAgendaDiaria(&ag)
	}

	catalogo := make([]response_catalogo.CatalogoResponse, len(p.Catalogo))
	for i, c := range p.Catalogo {
		catalogo[i] = response_catalogo.FromCatalogo(&c)
	}

	return PrestadorResponse{
		ID:       p.ID,
		Nome:     p.Nome,
		Email:    p.Email,
		Telefone: p.Telefone,
		Ativo:    p.Ativo,
		Catalogo: catalogo,
		Agenda:   agenda,
	}
}
