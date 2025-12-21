package output

import (
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendamentoOutput struct {
	ID             string
	PrestadorNome  string
	PrestadorTel   string
	ServicoNome    string
	ServicoDuracao int
	ServicoPreco   int
	DataHoraInicio time.Time
	DataHoraFim    time.Time
	Status         domain.StatusDoAgendamento
	Notas          string
}
