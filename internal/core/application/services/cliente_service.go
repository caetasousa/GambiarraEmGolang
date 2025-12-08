package services

import (
	"errors"
	"meu-servico-agenda/internal/core/application/ports"
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
)

type CadastrarClienteInput struct {
	Nome     string `json:"nome" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Telefone string `json:"telefone"`
}

type CadastroDeCliente struct {
	repo ports.ClienteRepositorio
}

func NovoCadastradoDeCliente(r ports.ClienteRepositorio) *CadastroDeCliente {
	return &CadastroDeCliente{repo: r}
}

func (s *CadastroDeCliente) Executar(input CadastrarClienteInput) (*domain.Cliente, error) {

	novoCliente := &domain.Cliente{
		ID:       xid.New().String(), 
		Nome:     input.Nome,
		Email:    input.Email,
		Telefone: input.Telefone,
	}

	if err := s.repo.Salvar(novoCliente); err != nil {
		return nil, errors.New("falha ao salvar cliente: " + err.Error())
	}

	return novoCliente, nil
}