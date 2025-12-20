package input

import "time"

type CadastrarAgendamentoInput struct {
	ClienteID      string
	PrestadorID    string
	CatalogoID     string
	DataHoraInicio time.Time
	Notas          string
}
