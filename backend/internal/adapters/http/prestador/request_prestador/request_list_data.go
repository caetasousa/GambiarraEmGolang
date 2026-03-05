package request_prestador

import "meu-servico-agenda/internal/core/application/input"

type BuscarPrestadoresDataRequest struct {
	Data  string `form:"data" binding:"required,datetime=2006-01-02"`
	Page  int    `form:"page" binding:"omitempty,min=1"`          // Página (padrão: 1)
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"` // Limite por página (padrão: 10, máx: 100)
}

func (r *BuscarPrestadoresDataRequest) ToInputPrestador() *input.PrestadorListDataInput {
	// Validações e valores padrão
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

	return &input.PrestadorListDataInput{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
		Data:   r.Data,
	}
}
