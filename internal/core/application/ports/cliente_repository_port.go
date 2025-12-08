package ports

import "meu-servico-agenda/internal/core/domain"

type ClienteRepositorio interface {
    Salvar(cliente *domain.Cliente) error
}