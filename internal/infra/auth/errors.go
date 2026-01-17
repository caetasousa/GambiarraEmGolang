package auth

import "errors"

// Erros relacionados à infraestrutura de autenticação JWT
var (
	ErrTokenInvalido        = errors.New("token inválido")
	ErrTokenExpirado        = errors.New("token expirado")
	ErrTokenMalformed       = errors.New("token mal formatado")
	ErrCredenciaisInvalidas = errors.New("credenciais inválidas")
	ErrFalhaAoGerarToken    = errors.New("falha ao gerar token de autenticação")
)
