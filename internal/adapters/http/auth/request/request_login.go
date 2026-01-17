package request

import (
	"meu-servico-agenda/internal/core/application/input"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"admin@admin.com"`
	Senha string `json:"senha" binding:"required" example:"admin123"`
}

func (r *LoginRequest) ToLoginInput() *input.LoginInput {
	return &input.LoginInput{
		Email: r.Email,
		Senha: r.Senha,
	}
}

// ValidarSenha compara a senha fornecida com o hash armazenado
func ValidarSenha(passwordHash, senha string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(senha))
}
