// internal/adapter/mapper/prestador_mapper.go
package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

// PrestadorFromDomainOutput converte um único prestador de domain para output
func PrestadorFromDomainOutput(p *domain.Prestador) *output.BuscarPrestadorOutput {
	if p == nil {
		return nil
	}

	return &output.BuscarPrestadorOutput{
		ID:        p.ID,
		Nome:      p.Nome,
		Cpf:       p.Cpf,
		Email:     p.Email,
		Telefone:  p.Telefone,
		Ativo:     p.Ativo,
		ImagemUrl: p.ImagemUrl,
		Catalogo:  CatalogosFromDomain(p.Catalogo),
		Agenda:    AgendasFromDomain(p.Agenda),
	}
}

// PrestadoresFromDomainOutput converte lista de prestadores de domain para output
func PrestadoresFromDomainOutput(prestadores []*domain.Prestador) []*output.BuscarPrestadorOutput {
	if prestadores == nil {
		return []*output.BuscarPrestadorOutput{}
	}

	outputs := make([]*output.BuscarPrestadorOutput, 0, len(prestadores))
	for _, p := range prestadores {
		outputs = append(outputs, PrestadorFromDomainOutput(p))
	}

	return outputs
}

// CatalogosFromDomain converte catálogos de domain para output
func CatalogosFromDomainAll(catalogos []domain.Catalogo) []output.CatalogoOutput {
	if catalogos == nil {
		return []output.CatalogoOutput{}
	}

	outputs := make([]output.CatalogoOutput, 0, len(catalogos))
	for _, c := range catalogos {
		outputs = append(outputs, output.CatalogoOutput{
			ID:            c.ID,
			Nome:          c.Nome,
			DuracaoPadrao: c.DuracaoPadrao,
			Preco:         c.Preco,
			ImagemUrl:     c.ImagemUrl,
			Categoria:     c.Categoria,
		})
	}

	return outputs
}

// AgendasFromDomain converte agendas de domain para output
func AgendasFromDomain(agendas []domain.AgendaDiaria) []output.AgendaDiariaOutput {
	if agendas == nil {
		return []output.AgendaDiariaOutput{}
	}

	outputs := make([]output.AgendaDiariaOutput, 0, len(agendas))
	for _, a := range agendas {
		outputs = append(outputs, output.AgendaDiariaOutput{
			ID:         a.Id,
			Data:       a.Data,
			Intervalos: IntervalosFromDomain(a.Intervalos),
		})
	}

	return outputs
}

// IntervalosFromDomain converte intervalos de domain para output
func IntervalosFromDomain(intervalos []domain.IntervaloDiario) []output.IntervaloDiarioOutput {
	if intervalos == nil {
		return []output.IntervaloDiarioOutput{}
	}

	outputs := make([]output.IntervaloDiarioOutput, 0, len(intervalos))
	for _, i := range intervalos {
		outputs = append(outputs, output.IntervaloDiarioOutput{
			ID:         i.Id,
			// ✅ Converte time.Time para string no formato HH:MM:SS
			HoraInicio: i.HoraInicio.Format("15:04:05"),
			HoraFim:    i.HoraFim.Format("15:04:05"),
		})
	}

	return outputs
}