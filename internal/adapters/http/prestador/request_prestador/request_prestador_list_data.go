package request_prestador

import "meu-servico-agenda/internal/core/application/input"

type BuscarPrestadoresDataRequest struct {
	Data  string `form:"data" binding:"required,datetime=2006-01-02"`
	Page  int    `form:"page" binding:"omitempty,min=1"`          // Página (padrão: 1)
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"` // Limite por página (padrão: 10, máx: 100)
}

func (r *BuscarPrestadoresDataRequest) ToInputPrestador() *input.PrestadorListDataInput {
	input := &input.PrestadorListDataInput{
		Page:  r.Page,
		Limit: r.Limit,
		Data: *&r.Data, // ✅ Sempre vai ter valor
	}

	// Validações de paginação
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	return input
}
