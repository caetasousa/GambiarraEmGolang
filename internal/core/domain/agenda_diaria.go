package domain

import (
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

func NovaAgendaDiaria(data time.Time, intervalos []IntervaloDiario) (*AgendaDiaria, error) {
	if len(intervalos) == 0 {
		return nil, ErrAgendaSemIntervalos
	}

	for _, it := range intervalos {
		if !it.HoraInicio.Before(it.HoraFim) {
			return nil, ErrIntervaloHorarioInvalido
		}
	}
	err := ValidarDataNoPassado(data)
	if err != nil {
		return nil, err
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

func (a *AgendaDiaria) PermiteAgendamento(inicio, fim time.Time) bool {

	// Converte entrada para UTC
	inicioUTC := inicio.UTC()
	fimUTC := fim.UTC()

	// Parse da data da agenda EM UTC
	dataAgenda, err := time.ParseInLocation("2006-01-02", a.Data, time.UTC)
	if err != nil {
		return false
	}

	for _, it := range a.Intervalos {

		inicioIntervalo := time.Date(
			dataAgenda.Year(),
			dataAgenda.Month(),
			dataAgenda.Day(),
			it.HoraInicio.Hour(),
			it.HoraInicio.Minute(),
			0,
			0,
			time.UTC,
		)

		fimIntervalo := time.Date(
			dataAgenda.Year(),
			dataAgenda.Month(),
			dataAgenda.Day(),
			it.HoraFim.Hour(),
			it.HoraFim.Minute(),
			0,
			0,
			time.UTC,
		)

		if !inicioUTC.Before(inicioIntervalo) && !fimUTC.After(fimIntervalo) {
			//log.Printf("✅ retornei true, %s data inicio e %s data final", inicioIntervalo, fimIntervalo)
			return true
		}
	}
	//log.Printf("✅ retornei false")
	return false
}
