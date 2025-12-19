package input

import "time"

type AdicionarAgendaInput struct {
	Data       time.Time
	Intervalos []IntervaloInput
}

type IntervaloInput struct {
	Inicio time.Time
	Fim    time.Time
}
