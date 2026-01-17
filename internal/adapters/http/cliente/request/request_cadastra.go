package request

import (
	"meu-servico-agenda/internal/core/application/input"

	"golang.org/x/crypto/bcrypt"
)

type ClienteRequest struct {
	Nome     string `json:"nome" binding:"required,min=3,max=100" example:"João da Silva"`
	Email    string `json:"email" binding:"required,email" example:"joao@email.com"`
	Telefone string `json:"telefone" binding:"required,min=8,max=15" example:"62999677481"`
	Senha    string `json:"senha" binding:"required,min=8" example:"senha123"`
}

func (r *ClienteRequest) ToCadastrarClienteInput() (*input.CadastrarClienteInput, error) {
	// Gerar hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &input.CadastrarClienteInput{
		Nome:     r.Nome,
		Email:    r.Email,
		Telefone: r.Telefone,
		Senha:    string(hash), // Passa o hash ao invés da senha em texto
	}, nil
}
