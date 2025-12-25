package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

func FromDomainToCriarOutput(p *domain.Prestador) *output.CriarPrestadorOutput {
	return &output.CriarPrestadorOutput{
		ID:       p.ID,
		Nome:     p.Nome,
		Email:    p.Email,
		Telefone: p.Telefone,
		Ativo:    p.Ativo,
		Catalogo: CatalogosFromDomain(p.Catalogo),
	}
}

func CatalogosFromDomain(catalogos []domain.Catalogo) []output.CatalogoOutput {
	result := make([]output.CatalogoOutput, len(catalogos))
	for i, c := range catalogos {
		result[i] = output.CatalogoOutput{
			ID:            c.ID,
			Nome:          c.Nome,
			DuracaoPadrao: c.DuracaoPadrao,
			Preco:         c.Preco,
			Categoria:     c.Categoria,
			ImagemUrl:     c.ImagemUrl,
		}
	}
	return result
}
