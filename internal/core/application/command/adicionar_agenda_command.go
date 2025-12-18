package command

import "time"

type AdicionarAgendaCommand struct {
	Data       time.Time
	Intervalos []IntervaloCommand
}

type IntervaloCommand struct {
	Inicio time.Time
	Fim    time.Time
}
