package domain

import (
	"errors"
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
	ClienteID      string
	PrestadorID    string
	CatalogoID     string
	DataHoraInicio time.Time
	DataHoraFim    time.Time
	Status         StatusDoAgendamento
	Notas          string
}

func NovoAgendamento(
	clienteID string,
	prestadorID string,
	catalogoID string,
	dataHoraInicio time.Time,
	dataHoraFim time.Time,
	nota string,
) (*Agendamento, error) {

	if !dataHoraInicio.Before(dataHoraFim) {
		return nil, errors.New("horário início deve ser antes do fim")
	}

	return &Agendamento{
		ID:             xid.New().String(),
		ClienteID:      clienteID,
		PrestadorID:    prestadorID,
		CatalogoID:     catalogoID,
		DataHoraInicio: dataHoraInicio,
		DataHoraFim:    dataHoraFim,
		Status:         Pendente,
		Notas:          nota,
	}, nil
}
