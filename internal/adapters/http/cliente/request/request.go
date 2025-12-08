package request

import (
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
)

type ClienteRequest struct {
	Nome     string `json:"nome" binding:"required"`
	Email    string `json:"email" binding:"email"`
	Telefone string `json:"telefone" binding:"required"`
}

func (r *ClienteRequest) ToCliente() *domain.Cliente {
	return &domain.Cliente{
		ID:       xid.New().String(),
		Nome:     r.Nome,
		Email:    r.Email,
		Telefone: r.Telefone,
	}
}
