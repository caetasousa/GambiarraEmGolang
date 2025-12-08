package domain

import "time"

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
	ServicoID      string
	DataHoraInicio time.Time
	DataHoraFim    time.Time // Calculado com base na duração do Serviço
	Status         StatusDoAgendamento    // Ex: "Confirmado", "Cancelado", "Pendente"
	Notas          string
}
