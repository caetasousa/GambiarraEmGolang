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

func (s *ServiceCliente) Cadastra(input domain.Cliente) (*domain.Cliente, error) {

	if err := s.repo.Salvar(&input); err != nil {
		return nil, errors.New("falha ao salvar cliente: " + err.Error())
	}

	return &input, nil
}

func (s *ServiceCliente) BuscarPorId(id string) (*domain.Cliente, error) {
	cliente, err := s.repo.BuscarPorId(id)

	// 1. Erro de Infraestrutura (DB offline, etc.)
	if err != nil {
		return nil, errors.New("falha na infraestrutura ao buscar cliente")
	}

	// 2. Cliente Não Encontrado (Nil retornado pelo repositório)
	if cliente == nil {
		return nil, errors.New("cliente não encontrado") // Erro de negócio para o Controller/Handler
	}

	// 3. Sucesso
	return cliente, nil
}
