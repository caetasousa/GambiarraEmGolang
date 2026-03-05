package domain

type Cliente struct {
	ID           string
	Nome         string
	Email        string
	Telefone     string
	PasswordHash string // Senha hasheada com bcrypt
	Role         Role   // Sempre RoleCliente
}

// NovoCliente cria um novo cliente recebendo o id e o hash da senha já pronto
func NovoCliente(id, nome, email, telefone, passwordHash string) *Cliente {
	return &Cliente{
		ID:           id,
		Nome:         nome,
		Email:        email,
		Telefone:     telefone,
		PasswordHash: passwordHash, // Recebe o hash já pronto
		Role:         RoleCliente,
	}
}
