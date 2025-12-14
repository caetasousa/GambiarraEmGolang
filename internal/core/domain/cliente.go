package domain

import (
	"github.com/rs/xid"
)

type Cliente struct {
	ID       string
	Nome     string
	Email    string
	Telefone string
}

func NovoCliente(nome, email, telefone string) (*Cliente, error) {
	return &Cliente{
		ID:       xid.New().String(),
		Nome:     nome,
		Email:    email,
		Telefone: telefone,
	}, nil
}
