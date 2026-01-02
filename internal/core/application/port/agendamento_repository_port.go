package port

import (
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendamentoRepositorio interface {
	CriaAgendamento(*domain.Agendamento) error
	BuscarPorPrestadorEPeriodo(prestadorID string, inicio time.Time, fim time.Time) ([]*domain.Agendamento, error)
	BuscarPorClienteEPeriodo(clienteID string, inicio time.Time, fim time.Time) ([]*domain.Agendamento, error)
	BuscarAgendamentoClienteAPartirDaData(clienteID string, data time.Time) ([]*domain.Agendamento, error)
	BuscarAgendamentoPrestadorAPartirDaData(clienteID string, data time.Time) ([]*domain.Agendamento, error)
}
