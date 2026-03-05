package output

// ============================================================================
// Estruturas de Dados de Usuário para Login
// ============================================================================

// LoginClienteData representa os dados de um cliente após login (SEM senha)
type LoginClienteData struct {
	ID       string `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

// LoginPrestadorData representa os dados de um prestador após login (SEM senha)
type LoginPrestadorData struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Telefone  string `json:"telefone"`
	Cpf       string `json:"cpf"`
	ImagemUrl string `json:"imagem_url"`
}

// ============================================================================
// Estruturas de Resposta de Login
// ============================================================================

// LoginOutput representa a resposta unificada após login bem-sucedido
type LoginOutput struct {
	Token string      `json:"token"`
	Role  string      `json:"role"`
	User  interface{} `json:"user"`
}


