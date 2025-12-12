package port

import "meu-servico-agenda/internal/core/domain"

type AgendaDiariaRepositorio interface {
	Salvar(agoenda *domain.AgendaDiaria) error
}
