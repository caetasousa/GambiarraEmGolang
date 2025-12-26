package input

type CadastrarPrestadorInput struct {
	Nome        string
	CPF         string
	Email       string
	Telefone    string
	ImagemUrl   string
	CatalogoIDs []string
}
