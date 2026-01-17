package domain

// Role representa os papéis de usuário no sistema
type Role string

const (
	RoleAdmin     Role = "admin"     // Acesso total ao sistema
	RoleCliente   Role = "cliente"   // Cliente que agenda serviços
	RolePrestador Role = "prestador" // Prestador de serviços
)

// IsValid verifica se a role é válida
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleCliente, RolePrestador:
		return true
	default:
		return false
	}
}

// String retorna a representação em string da role
func (r Role) String() string {
	return string(r)
}
