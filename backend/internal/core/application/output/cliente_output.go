package output

// BuscarClienteOutput representa os dados de um cliente (SEM senha)
type BuscarClienteOutput struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
	Role     string `json:"role"`
}
