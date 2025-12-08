package repository

import (
	"meu-servico-agenda/internal/core/domain"
)

// FakeClienteRepositorio é uma implementação de ClienteRepositorio que armazena dados em memória.
type FakeClienteRepositorio struct {
	Clientes map[string]*domain.Cliente 
}

// NewFakeClienteRepositorio cria e inicializa o repositório fake.
func NewFakeClienteRepositorio() *FakeClienteRepositorio {
	return &FakeClienteRepositorio{
		Clientes: make(map[string]*domain.Cliente),
	}
}

// Salvar simula a persistência.
func (r *FakeClienteRepositorio) Salvar(cliente *domain.Cliente) error {
	r.Clientes[cliente.Email] = cliente
	return nil
}

// BuscarPorEmail simula a busca por um email.
func (r *FakeClienteRepositorio) BuscarPorEmail(email string) (*domain.Cliente, error) {
	cliente, ok := r.Clientes[email]
	if !ok {
		// Retornar nil, nil para "não encontrado" é um padrão comum em repositórios Go.
		return nil, nil 
	}
	return cliente, nil
}