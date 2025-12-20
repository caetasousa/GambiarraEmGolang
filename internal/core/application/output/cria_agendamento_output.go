package output

import "time"

type AgendamentoOutput struct {
	ID            string
	PrestadorNome string
	PrestadorTel  string
	ServicoNome   string
	ServicoDuracao int
	ServicoPreco   int
	DataHoraInicio time.Time
	DataHoraFim    time.Time
	Status         int
	Notas          string
}
