package repository

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"

	"github.com/klassmann/cpfcnpj"
)

type FakePrestadorRepositorio struct {
	storage map[string]*domain.Prestador
}

func NovoFakePrestadorRepositorio() port.PrestadorRepositorio {
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

func (r *FakePrestadorRepositorio) BuscarPorCPF(cpf string) (*domain.Prestador, error) {
	cpf = cpfcnpj.Clean(cpf)
	for _, p := range r.storage {
		if cpfcnpj.Clean(p.Cpf) == cpf {
			return p, nil
		}
	}
	return nil, nil
}

func (r *FakePrestadorRepositorio) BuscarAgendaDoDia(prestadorID string, data string) (*domain.AgendaDiaria, error) {
	prestador, ok := r.storage[prestadorID]
	if !ok {
		return nil, nil
	}

	for _, agenda := range prestador.Agenda {
		if agenda.Data == data {
			return &agenda, nil
		}
	}

	return nil, nil
}