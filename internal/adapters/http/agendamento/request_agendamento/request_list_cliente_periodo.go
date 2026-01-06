package request_agendamento

import (
	"meu-servico-agenda/internal/adapters/http"
	"meu-servico-agenda/internal/core/application/input"
	"time"
)

type AgendamentoClientePeriodoRequest struct {
	DataInicio string `form:"data_inicio" binding:"required,datetime=2006-01-02" example:"2025-01-01"`
	DataFim    string `form:"data_fim" binding:"required,datetime=2006-01-02" example:"2025-01-31"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
}

func (req *AgendamentoClientePeriodoRequest) ToInput() (*input.ListarAgendamentoClientePeriodoInput, error) {
	dataInicio, err := time.Parse("2006-01-02", req.DataInicio)
	if err != nil {
		return nil, http.ErrFormatoDataHoraInvalido
	}

	dataFim, err := time.Parse("2006-01-02", req.DataFim)
	if err != nil {
		return nil, http.ErrFormatoDataHoraInvalido
	}

	return &input.ListarAgendamentoClientePeriodoInput{
		DataInicio: dataInicio,
		DataFim:    dataFim,
		Page:       req.Page,
		Limit:      req.Limit,
	}, nil
}
