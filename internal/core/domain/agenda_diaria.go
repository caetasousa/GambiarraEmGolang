package domain

import (
	"sort"
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

	// Validar cada intervalo individualmente
	for _, it := range intervalos {
		if !it.HoraInicio.Before(it.HoraFim) {
			return nil, ErrIntervaloHorarioInvalido
		}
	}

	// ✅ NOVO: Validar sobreposição de intervalos
	if err := validarSobreposicao(intervalos); err != nil {
		return nil, err
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

func validarSobreposicao(intervalos []IntervaloDiario) error {
	if len(intervalos) < 2 {
		return nil // Nada a validar
	}

	// Criar cópia para não modificar o slice original
	ordenados := make([]IntervaloDiario, len(intervalos))
	copy(ordenados, intervalos)

	// Ordenar por hora de início
	sort.Slice(ordenados, func(i, j int) bool {
		hi := ordenados[i].HoraInicio
		hj := ordenados[j].HoraInicio
		return hi.Hour() < hj.Hour() || (hi.Hour() == hj.Hour() && hi.Minute() < hj.Minute())
	})

	// Verificar sobreposição
	for i := 0; i < len(ordenados)-1; i++ {
		fimAtual := ordenados[i].HoraFim
		inicioProximo := ordenados[i+1].HoraInicio

		// Fim do atual deve ser ANTES do início do próximo
		// Se fimAtual >= inicioProximo, há sobreposição
		if fimAtual.Hour() > inicioProximo.Hour() ||
			(fimAtual.Hour() == inicioProximo.Hour() && fimAtual.Minute() >= inicioProximo.Minute()) {
			return ErrIntervalosSesobrepoe
		}
	}

	return nil
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
