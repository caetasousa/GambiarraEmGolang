package input

type CadastrarPrestadorInput struct {
	ID          string
	Nome        string
	CPF         string
	Email       string
	Telefone    string
	Senha       string
	ImagemUrl   string
	CatalogoIDs []string
}
