package request_agendamento

import (
	"meu-servico-agenda/internal/adapters/http"
	"meu-servico-agenda/internal/core/application/input"
	"time"
)

type AgendamentoPrestadorPeriodoRequest struct {
	DataInicio string `form:"data_inicio" binding:"required,datetime=2006-01-02" example:"2025-01-01"`
	DataFim    string `form:"data_fim" binding:"required,datetime=2006-01-02" example:"2025-01-31"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

func (req *AgendamentoPrestadorPeriodoRequest) ToInput() (*input.ListarAgendamentoPrestadorPeriodoInput, error) {
	dataInicio, err := time.Parse("2006-01-02", req.DataInicio)
	if err != nil {
		return nil, http.ErrFormatoDataHoraInvalido
	}

	dataFim, err := time.Parse("2006-01-02", req.DataFim)
	if err != nil {
		return nil, http.ErrFormatoDataHoraInvalido
	}

	// Validações de paginação
	page := req.Page
	limit := req.Limit

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return &input.ListarAgendamentoPrestadorPeriodoInput{
		DataInicio: dataInicio,
		DataFim:    dataFim,
		Page:       page,
		Limit:      limit,
		Offset:     (page - 1) * limit,
	}, nil
}
