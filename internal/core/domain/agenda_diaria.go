package domain

import (
	"errors"
	"time"

	"github.com/rs/xid"
)

type AgendaDiaria struct {
	Id         string
	Data       string            // Ex: "2026-01-15" (A data em que o trabalho ocorre)
	Intervalos []IntervaloDiario // Lista de horários disponíveis naquele dia
}

// IntervaloDiario: O bloco de tempo (Ex: 09:00 - 12:00)
type IntervaloDiario struct {
	Id         string
	HoraInicio time.Time
	HoraFim    time.Time
}

var (
	ErrAgendaSemIntervalos      = errors.New("agenda deve conter ao menos um intervalo")
	ErrIntervaloHorarioInvalido = errors.New("hora início deve ser menor que hora fim")
)

func NovaAgendaDiaria(data time.Time, intervalos []IntervaloDiario) (*AgendaDiaria, error) {
	if len(intervalos) == 0 {
		return nil, ErrAgendaSemIntervalos
	}

	for _, it := range intervalos {
		if !it.HoraInicio.Before(it.HoraFim) {
			return nil, ErrIntervaloHorarioInvalido
		}
	}

	return &AgendaDiaria{
		Id:         xid.New().String(),
		Data:       data.Format("2006-01-02"),
		Intervalos: intervalos,
	}, nil
}

func NovoIntervaloDiario(horaInicio, horaFim time.Time) (*IntervaloDiario, error) {
	if !horaInicio.Before(horaFim) {
		return nil, ErrIntervaloHorarioInvalido
	}

	return &IntervaloDiario{
		Id:         xid.New().String(),
		HoraInicio: horaInicio,
		HoraFim:    horaFim,
	}, nil
}
