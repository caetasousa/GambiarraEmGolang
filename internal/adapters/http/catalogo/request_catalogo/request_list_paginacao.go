package request_catalogo

import "meu-servico-agenda/internal/core/application/input"

type CatalogoListRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (cat *CatalogoListRequest) ToInputCatalogo() *input.ListCatalogoInput {
	return &input.ListCatalogoInput{
		Page:  cat.Page,
		Limit: cat.Limit,
	}
}
