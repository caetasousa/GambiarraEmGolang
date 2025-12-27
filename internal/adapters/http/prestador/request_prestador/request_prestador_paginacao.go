package request_prestador

import "meu-servico-agenda/internal/core/application/input"

type PrestadorListRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (cat *PrestadorListRequest) ToInputPrestador() *input.PrestadorListInput {
	return &input.PrestadorListInput{
		Page:  cat.Page,
		Limit: cat.Limit,
	}
}
