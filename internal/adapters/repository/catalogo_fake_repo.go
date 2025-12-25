package repository

import (
	"errors"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"sort"
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

func (r *CatalogoFakeRepo) Listar(limit, offset int) ([]*domain.Catalogo, error) {
	catalogos := make([]*domain.Catalogo, 0, len(r.Catalogo))

	// map → slice
	for _, c := range r.Catalogo {
		catalogos = append(catalogos, c)
	}

	// ordenação previsível (importante para testes)
	sort.Slice(catalogos, func(i, j int) bool {
		return catalogos[i].ID < catalogos[j].ID
	})

	if offset >= len(catalogos) {
		return []*domain.Catalogo{}, nil
	}

	end := offset + limit
	if end > len(catalogos) {
		end = len(catalogos)
	}

	return catalogos[offset:end], nil
}

func (r *CatalogoFakeRepo) Contar() (int, error) {
	return len(r.Catalogo), nil
}

func (r *CatalogoFakeRepo) Atualizar(catalogo *domain.Catalogo) error {
	// Verifica se o catálogo existe
	if _, exists := r.Catalogo[catalogo.ID]; !exists {
		return errors.New("catálogo não encontrado")
	}
	
	// Atualiza o catálogo no map
	r.Catalogo[catalogo.ID] = catalogo
	return nil
}

func (r *CatalogoFakeRepo) Deletar(id string) error {
	// Verifica se o catálogo existe
	if _, exists := r.Catalogo[id]; !exists {
		return errors.New("catálogo não encontrado")
	}
	
	// Remove o catálogo do map
	delete(r.Catalogo, id)
	return nil
}