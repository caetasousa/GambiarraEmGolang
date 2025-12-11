package request

import (
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
)

type ClienteRequest struct {
	Nome     string `json:"nome" binding:"required,min=3,max=100" swagger:"desc('Nome do cliente')"`
	Email    string `json:"email" binding:"omitempty,email" swagger:"desc('Email do cliente')"`
	Telefone string `json:"telefone" binding:"required,min=8,max=15" swagger:"desc('Telefone do cliente')"`
}

func (r *ClienteRequest) ToCliente() *domain.Cliente {
	return &domain.Cliente{
		ID:       xid.New().String(),
		Nome:     r.Nome,
		Email:    r.Email,
		Telefone: r.Telefone,
	}
}
