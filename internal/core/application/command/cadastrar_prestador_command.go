package command

type CadastrarPrestadorCommand struct {
	Nome        string
	CPF         string
	Email       string
	Telefone    string
	CatalogoIDs []string
}
