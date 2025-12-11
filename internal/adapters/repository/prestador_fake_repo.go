package repository

import "meu-servico-agenda/internal/core/domain"

type FakePrestadorRepositorio struct {
	storage map[string]*domain.Prestador
}

func NovoFakePrestadorRepositorio() *FakePrestadorRepositorio {
	return &FakePrestadorRepositorio{
		storage: make(map[string]*domain.Prestador),
	}
}

func (r *FakePrestadorRepositorio) Salvar(prestador *domain.Prestador) error {
	r.storage[prestador.ID] = prestador
	return nil
}

func (r *FakePrestadorRepositorio) BuscarPorId(id string) (*domain.Prestador, error) {
	prestador, ok := r.storage[id]
	if !ok {
		return nil, nil
	}
	return prestador, nil
}
