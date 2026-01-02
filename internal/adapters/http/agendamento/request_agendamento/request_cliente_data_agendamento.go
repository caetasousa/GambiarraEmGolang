package request_agendamento

import (
	"fmt"
	"meu-servico-agenda/internal/core/application/input"
	"time"
)

type AgendamentoDataRequest struct {
    Data string `form:"data" binding:"required,datetime=2006-01-02" example:"2025-01-03"`
}

func (r *AgendamentoDataRequest) ToAgendamentoDataInput() (*input.AgendamentoDataInput, error) {
    data, err := time.Parse("2006-01-02", r.Data)
    if err != nil {
        return nil, fmt.Errorf("data inválida: %w", err)
    }

    return &input.AgendamentoDataInput{
        Data: data,
    }, nil
}
