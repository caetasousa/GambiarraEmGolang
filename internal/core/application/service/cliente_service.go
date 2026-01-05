package service

import (
	"errors"

	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type ServiceCliente struct {
	repo port.ClienteRepositorio
}

func NovoServiceCliente(r port.ClienteRepositorio) *ServiceCliente {
	return &ServiceCliente{repo: r}
}

func (s *ServiceCliente) Cadastra(cliente *domain.Cliente) (*domain.Cliente, error) {
	if cliente == nil {
		return nil, ErrClienteInvalido
	}

	if err := s.repo.Salvar(cliente); err != nil {
		if errors.Is(err, repository.ErrEmailDuplicado) {
			return nil, repository.ErrEmailDuplicado
		}
		return nil, ErrFalhaInfraestrutura
	}

	return cliente, nil
}

func (s *ServiceCliente) BuscarPorId(id string) (*domain.Cliente, error) {
	cliente, err := s.repo.BuscarPorId(id)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	// BuscarPorId retorna (nil, nil) se não encontrar
	if cliente == nil {
		return nil, ErrClienteNaoEncontrado
	}

	return cliente, nil
}
