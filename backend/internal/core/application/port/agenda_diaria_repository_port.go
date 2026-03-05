package port

import (
	"meu-servico-agenda/internal/core/domain"
)

type AgendaDiariaRepositorio interface {
	Salvar(agenda *domain.AgendaDiaria, prestadorId string) error
	AtualizarAgenda(agenda *domain.AgendaDiaria, prestadorID string) error 
	BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error)
	DeletarAgenda(prestadorID string, data string) error
}