package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

func CatalogoFromDomainOutput(c *domain.Catalogo) *output.CatalogoOutput {
	if c == nil {
		return nil
	}

	return &output.CatalogoOutput{
		ID:            c.ID,
		Nome:          c.Nome,
		DuracaoPadrao: c.DuracaoPadrao,
		Preco:         c.Preco,
		Categoria:     c.Categoria,
		ImagemUrl:     c.ImagemUrl, 
	}
}


func CatalogosFromDomainOutput(catalogos []*domain.Catalogo) []*output.CatalogoOutput {
	result := make([]*output.CatalogoOutput, 0, len(catalogos))

	for _, c := range catalogos {
		result = append(result, CatalogoFromDomainOutput(c))
	}

	return result
}