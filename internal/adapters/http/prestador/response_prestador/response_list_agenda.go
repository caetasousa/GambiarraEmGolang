package response_prestador

import (
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/core/application/output"
)

type PrestadorResponse struct {
	ID        string                               `json:"id"`
	Nome      string                               `json:"nome"`
	Email     string                               `json:"email"`
	Telefone  string                               `json:"telefone"`
	Cpf       string                               `json:"cpf"`
	Ativo     bool                                 `json:"ativo"`
	ImagemUrl string                               `json:"image_url"`
	Catalogo  []response_catalogo.CatalogoResponse `json:"catalogo"`
	Agenda    []AgendaDiariaResponse               `json:"agenda"`
}

func FromPrestadorOutput(o output.BuscarPrestadorOutput) PrestadorResponse {

	agenda := make([]AgendaDiariaResponse, len(o.Agenda))
	for i, ag := range o.Agenda {
		agenda[i] = FromAgendaDiariaOutput(ag)
	}

	catalogo := make([]response_catalogo.CatalogoResponse, len(o.Catalogo))
	for i, c := range o.Catalogo {
		catalogo[i] = response_catalogo.FromCatalogoResponse(c)
	}

	return PrestadorResponse{
		ID:        o.ID,
		Nome:      o.Nome,
		Email:     o.Email,
		Cpf:       o.Cpf,
		Telefone:  o.Telefone,
		Ativo:     o.Ativo,
		ImagemUrl: o.ImagemUrl,
		Catalogo:  catalogo,
		Agenda:    agenda,
	}
}
