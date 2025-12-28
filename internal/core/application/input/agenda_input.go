package input

import "time"

type AdicionarAgendaInput struct {
	PrestadorID string
	Data       time.Time
	Intervalos []IntervaloInput
}

type AtualizarAgendaInput struct {
	PrestadorID string
	Data        time.Time
	Intervalos  []IntervaloInput
}

type IntervaloInput struct {
	Inicio time.Time
	Fim    time.Time
}
