package input

type AlterarPrestadorInput struct {
	Id          string
	Nome        string
	Email       string
	Telefone    string
	ImagemUrl   string
	CatalogoIDs []string
}
