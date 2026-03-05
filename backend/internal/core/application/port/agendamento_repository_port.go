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
	ListarPorPrestadorEPeriodoPaginado(prestadorID string, inicio time.Time, fim time.Time, limit, offset int) ([]*domain.Agendamento, error)
	ListarPorClienteEPeriodoPaginado(clienteID string, inicio time.Time, fim time.Time, limit, offset int) ([]*domain.Agendamento, error)
	ContarPorPrestadorEPeriodo(prestadorID string, inicio time.Time, fim time.Time) (int, error)
	ContarPorClienteEPeriodo(clienteID string, inicio time.Time, fim time.Time) (int, error)
}
