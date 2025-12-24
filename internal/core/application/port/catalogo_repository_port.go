package port

import "meu-servico-agenda/internal/core/domain"

type CatalogoRepositorio interface {
	Salvar(catalogo *domain.Catalogo) error
	BuscarPorId(id string) (*domain.Catalogo, error)
	// Paginação
	Listar(limit, offset int) ([]*domain.Catalogo, error)
	Contar() (int, error)
}