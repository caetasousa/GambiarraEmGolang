package request

import (
	"meu-servico-agenda/internal/core/domain"
)

type ClienteRequest struct {
	Nome     string `json:"nome" binding:"required,min=3,max=100" example:"João da Silva" swagger:"desc('Nome do cliente')"`
	Email    string `json:"email" binding:"omitempty,email" example:"joao@email.com" swagger:"desc('Email do cliente')"`
	Telefone string `json:"telefone" binding:"required,min=8,max=15" example:"62999677481" swagger:"desc('Telefone do cliente')"`
}

func (r *ClienteRequest) ToCliente() (*domain.Cliente, error) {
	return domain.NovoCliente(r.Nome, r.Email, r.Telefone)
}
