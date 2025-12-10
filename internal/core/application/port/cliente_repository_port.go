package port

import "meu-servico-agenda/internal/core/domain"

type ClienteRepositorio interface {
    Salvar(cliente *domain.Cliente) error
    BuscarPorId(id string) (*domain.Cliente, error)
}