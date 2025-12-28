package request_prestador

import "meu-servico-agenda/internal/core/application/input"

type PrestadorListRequest struct {
	Page  int  `form:"page" binding:"omitempty,min=1"`
	Limit int  `form:"limit" binding:"omitempty,min=1,max=100"`
	Ativo *bool `form:"ativo" binding:"required"`
}

func (r *PrestadorListRequest) ToInputPrestador() *input.PrestadorListInput {
	input := &input.PrestadorListInput{
		Page:  r.Page,
		Limit: r.Limit,
		Ativo: *r.Ativo, // ✅ Sempre vai ter valor
	}

	// Valores padrão
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	return input
}