package services

import (
	"errors"

	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/core/application/ports"
	"meu-servico-agenda/internal/core/domain"
)

type CadastroDeCliente struct {
	repo ports.ClienteRepositorio
}

func NovoCadastradoDeCliente(r ports.ClienteRepositorio) *CadastroDeCliente {
	return &CadastroDeCliente{repo: r}
}

func (s *CadastroDeCliente) Executar(input request.ClienteRequest) (*domain.Cliente, error) {
	novoCliente := input.ToCliente()

	if err := s.repo.Salvar(novoCliente); err != nil {
		return nil, errors.New("falha ao salvar cliente: " + err.Error())
	}

	return novoCliente, nil
}
