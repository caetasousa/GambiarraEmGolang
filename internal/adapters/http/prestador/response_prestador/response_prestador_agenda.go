package response_prestador

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/core/application/output"
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

func FromPrestadorOutput(o output.BuscarPrestadorOutput) PrestadorResponse {

	agenda := make([]AgendaDiariaResponse, len(o.Agenda))
	for i, ag := range o.Agenda {
		agenda[i] = FromAgendaDiariaOutput(ag)
	}

	catalogo := make([]response_catalogo.CatalogoResponse, len(o.Catalogo))
	for i, c := range o.Catalogo {
		catalogo[i] = response_catalogo.FromCatalogoOutput(c)
	}

	return PrestadorResponse{
		ID:       o.ID,
		Nome:     o.Nome,
		Email:    o.Email,
		Telefone: o.Telefone,
		Ativo:    o.Ativo,
		Catalogo: catalogo,
		Agenda:   agenda,
	}
}
