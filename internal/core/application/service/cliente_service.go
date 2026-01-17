package service

import (
	"errors"

	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/mapper"
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type ServiceCliente struct {
	repo          port.ClienteRepositorio
	prestadorRepo port.PrestadorRepositorio
}

func NovoServiceCliente(r port.ClienteRepositorio, prestadorRepo port.PrestadorRepositorio) *ServiceCliente {
	return &ServiceCliente{
		repo:          r,
		prestadorRepo: prestadorRepo,
	}
}

func (s *ServiceCliente) Cadastra(inputData *input.CadastrarClienteInput) (*output.BuscarClienteOutput, error) {
	// Verificar se email já existe na tabela de clientes
	existente, _ := s.repo.BuscarPorEmail(inputData.Email)
	if existente != nil {
		return nil, repository.ErrEmailDuplicado
	}

	// Verificar se email já está sendo usado por um prestador
	if s.prestadorRepo != nil {
		prestadorExistente, err := s.prestadorRepo.BuscarPorEmail(inputData.Email)
		if err != nil {
			return nil, ErrFalhaInfraestrutura
		}
		if prestadorExistente != nil {
			return nil, ErrEmailJaUsadoPorPrestador
		}
	}

	// Criar cliente (senha já vem hasheada do request)
	cliente := domain.NovoCliente(
		inputData.ID,
		inputData.Nome,
		inputData.Email,
		inputData.Telefone,
		inputData.Senha, // Senha já é o hash
	)

	// Salvar no banco
	if err := s.repo.Salvar(cliente); err != nil {
		if errors.Is(err, repository.ErrEmailDuplicado) {
			return nil, repository.ErrEmailDuplicado
		}
		return nil, ErrFalhaInfraestrutura
	}

	// Retornar output (SEM senha)
	outputData := mapper.ClienteToOutput(cliente)
	return &outputData, nil
}

func (s *ServiceCliente) BuscarPorId(id string) (*output.BuscarClienteOutput, error) {
	cliente, err := s.repo.BuscarPorId(id)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	// BuscarPorId retorna (nil, nil) se não encontrar
	if cliente == nil {
		return nil, ErrClienteNaoEncontrado
	}

	// Converter para output (SEM senha)
	outputData := mapper.ClienteToOutput(cliente)
	return &outputData, nil
}
