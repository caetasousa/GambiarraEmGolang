package domain

import (
	"time"

	"github.com/rs/xid"
)

type StatusDoAgendamento int

const (
	Pendente StatusDoAgendamento = iota + 1
	Confirmado
	Cancelado
	Concluido
)

type Agendamento struct {
	ID             string
	Cliente        *Cliente
	Prestador      *Prestador
	Catalogo       *Catalogo
	DataHoraInicio time.Time
	DataHoraFim    time.Time
	Status         StatusDoAgendamento
	Notas          string
}

func NovoAgendamento(
	cliente *Cliente,
	prestador *Prestador,
	catalogo *Catalogo,
	dataHoraInicio time.Time,
	dataHoraFim time.Time,
	nota string,
) (*Agendamento, error) {

	if !dataHoraInicio.Before(dataHoraFim) {
		return nil, ErrHoraInicialMenorQueFinal
	}

	return &Agendamento{
		ID:             xid.New().String(),
		Cliente:        cliente,
		Prestador:      prestador,
		Catalogo:       catalogo,
		DataHoraInicio: dataHoraInicio,
		DataHoraFim:    dataHoraFim,
		Status:         Pendente,
		Notas:          nota,
	}, nil
}
