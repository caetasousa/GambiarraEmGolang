package request_prestador

import (
	"fmt"
	"meu-servico-agenda/internal/core/application/input"

	"time"
)

type IntervaloDiarioRequest struct {
	HoraInicio string `json:"hora_inicio" example:"08:00" binding:"required,datetime=15:04"`
	HoraFim    string `json:"hora_fim" example:"12:00" binding:"required,datetime=15:04"`
}

type AgendaDiariaRequest struct {
	Data       string                   `json:"data" example:"2025-01-03" binding:"required,datetime=2006-01-02"`
	Intervalos []IntervaloDiarioRequest `json:"intervalos" binding:"required,dive"`
}

func (r *AgendaDiariaRequest) ToAdicionarAgendaInput() (*input.AdicionarAgendaInput, error) {
	data, err := time.Parse("2006-01-02", r.Data)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	intervalos, err := r.parseIntervalos()
	if err != nil {
		return nil, err
	}

	return &input.AdicionarAgendaInput{
		Data:       data,
		Intervalos: intervalos,
	}, nil
}

func (r *AgendaDiariaRequest) ToAtualizarAgendaInput() (*input.AtualizarAgendaInput, error) {
	data, err := time.Parse("2006-01-02", r.Data)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	intervalos, err := r.parseIntervalos()
	if err != nil {
		return nil, err
	}

	return &input.AtualizarAgendaInput{
		Data:       data,
		Intervalos: intervalos,
	}, nil
}

func (r *AgendaDiariaRequest) parseIntervalos() ([]input.IntervaloInput, error) {
	intervalos := make([]input.IntervaloInput, 0, len(r.Intervalos))
	
	for _, i := range r.Intervalos {
		inicio, err := time.Parse("15:04", i.HoraInicio)
		if err != nil {
			return nil, fmt.Errorf("hora_inicio inválida: %w", err)
		}

		fim, err := time.Parse("15:04", i.HoraFim)
		if err != nil {
			return nil, fmt.Errorf("hora_fim inválida: %w", err)
		}

		intervalos = append(intervalos, input.IntervaloInput{
			Inicio: inicio,
			Fim:    fim,
		})
	}

	return intervalos, nil
}