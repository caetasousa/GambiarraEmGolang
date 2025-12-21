package port

import (
	"meu-servico-agenda/internal/core/domain"
)

type AgendaDiariaRepositorio interface {
	Salvar(agenda *domain.AgendaDiaria, prestadorId string) error
}