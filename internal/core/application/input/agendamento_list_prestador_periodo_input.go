package input

import "time"

type ListarAgendamentoPrestadorPeriodoInput struct {
	DataInicio time.Time
	DataFim    time.Time
	Page       int
	Limit      int
	Offset     int
}
