package request_prestador

import "meu-servico-agenda/internal/core/application/input"

type PrestadorListRequest struct {
	Page  int   `form:"page" binding:"omitempty,min=1"`
	Limit int   `form:"limit" binding:"omitempty,min=1,max=100"`
	Ativo *bool `form:"ativo" binding:"required"`
}

func (r *PrestadorListRequest) ToInputPrestador() *input.PrestadorListInput {
	page := r.Page
	limit := r.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return &input.PrestadorListInput{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
		Ativo:  *r.Ativo,
	}
}