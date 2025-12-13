package domain

import (
	"errors"

	"github.com/rs/xid"
)

type Cliente struct {
	ID       string
	Nome     string
	Email    string
	Telefone string
}

func NovoCliente(nome, email, telefone string) (*Cliente, error) {
	if nome == "" {
		return nil, errors.New("nome do cliente não pode ser vazio")
	}

	return &Cliente{
		ID:       xid.New().String(), // ✅ ID gerado no domínio
		Nome:     nome,
		Email:    email,
		Telefone: telefone,
	}, nil
}
