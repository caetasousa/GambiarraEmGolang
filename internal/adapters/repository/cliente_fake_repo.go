package repository

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

// FakeClienteRepositorio é uma implementação de ClienteRepositorio que armazena dados em memória.
type FakeClienteRepositorio struct {
	// A chave do mapa DEVE ser o ID do cliente
	Clientes map[string]*domain.Cliente
}

// NewFakeClienteRepositorio cria e inicializa o repositório fake.
func NewFakeClienteRepositorio() port.ClienteRepositorio {
	return &FakeClienteRepositorio{
		Clientes: make(map[string]*domain.Cliente),
	}
}

// Salvar simula a persistência, usando o ID do cliente como chave.
func (r *FakeClienteRepositorio) Salvar(cliente *domain.Cliente) error {
	// CORREÇÃO: Usar o ID como chave (índice)
	r.Clientes[cliente.ID] = cliente
	return nil
}

// BuscarPorId simula a busca por um id.
func (r *FakeClienteRepositorio) BuscarPorId(id string) (*domain.Cliente, error) {
	// CORRETO: Agora a busca por 'id' funcionará, pois o mapa é indexado por ID.
	cliente, ok := r.Clientes[id]
	if !ok {
		return nil, nil
	}
	return cliente, nil
}
