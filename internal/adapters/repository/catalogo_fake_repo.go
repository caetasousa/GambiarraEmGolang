package repository

import (
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoFakeRepo struct {
	Catalogo map[string]*domain.Catalogo
}

func NovoCatalogoFakeRepo() *CatalogoFakeRepo {
	return &CatalogoFakeRepo{Catalogo: make(map[string]*domain.Catalogo)}
}

func (r *CatalogoFakeRepo) Salvar(catalogo *domain.Catalogo) error {
	r.Catalogo[catalogo.ID] = catalogo
	return nil
}

func (r *CatalogoFakeRepo) BuscarPorId(id string) (*domain.Catalogo, error) {
	catalogo, ok := r.Catalogo[id]
	if !ok {
		return nil, nil
	}
	return catalogo, nil
}
