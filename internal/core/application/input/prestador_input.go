package input

type CadastrarPrestadorInput struct {
	Nome        string
	CPF         string
	Email       string
	Telefone    string
	CatalogoIDs []string
}
