package request_catalogo

import "meu-servico-agenda/internal/core/application/input"

type CatalogoListRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (cat *CatalogoListRequest) ToInputCatalogo() *input.ListCatalogoInput {
	page := cat.Page
	limit := cat.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return &input.ListCatalogoInput{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}
