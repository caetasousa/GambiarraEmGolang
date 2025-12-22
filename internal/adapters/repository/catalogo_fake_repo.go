package repository

import (
	"errors"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoFakeRepo struct {
	Catalogo map[string]*domain.Catalogo
}

func NovoCatalogoFakeRepo() port.CatalogoRepositorio {
	return &CatalogoFakeRepo{Catalogo: make(map[string]*domain.Catalogo)}
}

func (r *CatalogoFakeRepo) Salvar(catalogo *domain.Catalogo) error {
	r.Catalogo[catalogo.ID] = catalogo
	return nil
}

func (r *CatalogoFakeRepo) BuscarPorId(id string) (*domain.Catalogo, error) {
	catalogo := r.Catalogo[id]
	if catalogo == nil {
		return nil, errors.New("não encontrado")
	}
	return catalogo, nil
}
