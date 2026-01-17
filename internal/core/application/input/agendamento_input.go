package input

import "time"

type CadastrarAgendamentoInput struct {
	ID             string
	ClienteID      string
	PrestadorID    string
	CatalogoID     string
	DataHoraInicio time.Time
	Notas          string
}
