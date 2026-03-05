package port

import "meu-servico-agenda/internal/core/domain"

type CatalogoRepositorio interface {
	//Post
	Salvar(catalogo *domain.Catalogo) error
	//GetById
	BuscarPorId(id string) (*domain.Catalogo, error)
	//GetAll Paginação
	Listar(limit, offset int) ([]*domain.Catalogo, error)
	Contar() (int, error)
	//Update
	Atualizar(catalogo *domain.Catalogo) error
	//Delete
	Deletar(id string) error
}