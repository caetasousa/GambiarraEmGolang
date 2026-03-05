package input

// CadastrarClienteInput representa os dados necessários para cadastrar um cliente
type CadastrarClienteInput struct {
	ID       string
	Nome     string
	Email    string
	Telefone string
	Senha    string
}

// LoginInput representa os dados necessários para fazer login
type LoginInput struct {
	Email string
	Senha string
}
