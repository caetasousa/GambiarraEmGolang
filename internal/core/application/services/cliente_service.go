package services

import (
	"errors"

	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/core/application/ports"
	"meu-servico-agenda/internal/core/domain"
)

type ServiceCliente struct {
	repo ports.ClienteRepositorio
}

func NovoServiceCliente(r ports.ClienteRepositorio) *ServiceCliente {
	return &ServiceCliente{repo: r}
}

func (s *ServiceCliente) Cadastra(input request.ClienteRequest) (*domain.Cliente, error) {
	novoCliente := input.ToCliente()

	if err := s.repo.Salvar(novoCliente); err != nil {
		return nil, errors.New("falha ao salvar cliente: " + err.Error())
	}

	return novoCliente, nil
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
