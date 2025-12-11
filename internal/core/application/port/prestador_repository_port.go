package port

import "meu-servico-agenda/internal/core/domain"

type PrestadorDeServicoRepositorio interface {
	Salvar(prestador *domain.Prestador) error
}
