package service

import (
	"errors"

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
		return nil, ErrClienteNaoEncontrado
	}

	if err := s.repo.Salvar(cliente); err != nil {
		return nil, errors.New(ErrAoSalvarCliente.Error() + err.Error())
	}

	return cliente, nil
}

func (s *ServiceCliente) BuscarPorId(id string) (*domain.Cliente, error) {
	cliente, err := s.repo.BuscarPorId(id)

	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	if cliente == nil {
		return nil, ErrClienteNaoEncontrado
	}

	return cliente, nil
}
