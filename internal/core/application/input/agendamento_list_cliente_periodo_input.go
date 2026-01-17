package input

import "time"

type ListarAgendamentoClientePeriodoInput struct {
	DataInicio time.Time
	DataFim    time.Time
	Page       int
	Limit      int
	Offset     int
}
