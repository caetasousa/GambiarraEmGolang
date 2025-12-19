package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

func FromPrestador(p *domain.Prestador) *output.BuscarPrestadorOutput {

	agenda := make([]output.AgendaDiariaOutput, len(p.Agenda))
	for i, ag := range p.Agenda {

		intervalos := make([]output.IntervaloDiarioOutput, len(ag.Intervalos))
		for j, it := range ag.Intervalos {
			intervalos[j] = output.IntervaloDiarioOutput{
				ID:         it.Id,
				HoraInicio: it.HoraInicio.Format("15:04"),
				HoraFim:    it.HoraFim.Format("15:04"),
			}
		}

		agenda[i] = output.AgendaDiariaOutput{
			ID:         ag.Id,
			Data:       ag.Data,
			Intervalos: intervalos,
		}
	}

	catalogo := make([]output.CatalogoOutput, len(p.Catalogo))
	for i, c := range p.Catalogo {
		catalogo[i] = output.CatalogoOutput{
			ID:            c.ID,
			Nome:          c.Nome,
			DuracaoPadrao: c.DuracaoPadrao,
			Preco:         c.Preco,
			Categoria:     c.Categoria,
		}
	}

	return &output.BuscarPrestadorOutput{
		ID:       p.ID,
		Nome:     p.Nome,
		Email:    p.Email,
		Telefone: p.Telefone,
		Ativo:    p.Ativo,
		Catalogo: catalogo,
		Agenda:   agenda,
	}
}
