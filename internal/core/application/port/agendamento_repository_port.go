package port

import "meu-servico-agenda/internal/core/domain"

type AgendamentoRepositorio interface {
	CriaAgendamento(*domain.Agendamento) error
}
