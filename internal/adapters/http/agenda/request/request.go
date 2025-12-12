package request

import (
	"fmt"
	"meu-servico-agenda/internal/core/domain"
	"time"

	"github.com/rs/xid"
)

type IntervaloDiarioRequest struct {
	HoraInicio string `json:"hora_inicio" example:"08:00" binding:"required,datetime=15:04"`
	HoraFim    string `json:"hora_fim" example:"12:00" binding:"required,datetime=15:04"`
}

type AgendaDiariaRequest struct {
	Data       string                   `json:"data" example:"2025-01-03" binding:"required,datetime=2006-01-02"`
	Intervalos []IntervaloDiarioRequest `json:"intervalos" binding:"required,dive"`
}

func (r *AgendaDiariaRequest) ToAgendaDiaria() (*domain.AgendaDiaria, error) {
	data, err := time.Parse("2006-01-02", r.Data)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	intervalos := make([]domain.IntervaloDiario, 0, len(r.Intervalos))
	for _, i := range r.Intervalos {
		inicio, err := time.Parse("15:04", i.HoraInicio)
		if err != nil {
			return nil, fmt.Errorf("hora_inicio inválida: %w", err)
		}

		fim, err := time.Parse("15:04", i.HoraFim)
		if err != nil {
			return nil, fmt.Errorf("hora_fim inválida: %w", err)
		}

		if !inicio.Before(fim) {
			return nil, fmt.Errorf("hora_inicio '%s' deve ser menor que hora_fim '%s'", i.HoraInicio, i.HoraFim)
		}

		intervalos = append(intervalos, domain.IntervaloDiario{
			Id:         xid.New().String(),
			HoraInicio: inicio,
			HoraFim:    fim,
		})
	}

	return domain.NovaAgendaDiaria(data, intervalos)
}
