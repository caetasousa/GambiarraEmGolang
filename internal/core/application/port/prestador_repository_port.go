package port

import "meu-servico-agenda/internal/core/domain"

type PrestadorRepositorio interface {
	Salvar(prestador *domain.Prestador) error
	BuscarPorId(id string) (*domain.Prestador, error)
}
