package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

// AgendaDiaria: A estrutura central que define a agenda para uma data específica.
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

func NovaAgendaDiaria(data time.Time, intervalos []IntervaloDiario) (*AgendaDiaria, error) {
	if len(intervalos) == 0 {
		return nil, errors.New("a agenda deve conter ao menos um intervalo")
	}

	// validar cada intervalo
	for _, it := range intervalos {
		if !it.HoraInicio.Before(it.HoraFim) {
			return nil, fmt.Errorf("intervalo inválido: hora de início (%s) deve ser anterior à hora de fim (%s)",
				it.HoraInicio.Format("15:04"), it.HoraFim.Format("15:04"))
		}
	}

	return &AgendaDiaria{
		Id:         xid.New().String(),
		Data:       data.Format("2006-01-02"), // porque vc usa string no domínio
		Intervalos: intervalos,
	}, nil
}

func NovoIntervaloDiario(horaInicio, horaFim time.Time) (*IntervaloDiario, error) {
	if !horaInicio.Before(horaFim) {
		return nil, fmt.Errorf("hora_inicio %s deve ser menor que hora_fim %s", horaInicio.Format("15:04"), horaFim.Format("15:04"))
	}

	return &IntervaloDiario{
		Id:         xid.New().String(),
		HoraInicio: horaInicio,
		HoraFim:    horaFim,
	}, nil
}