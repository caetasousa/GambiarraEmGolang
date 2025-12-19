package output

type CriarPrestadorOutput struct {
	ID       string
	Nome     string
	Email    string
	Telefone string
	Ativo    bool
	Catalogo []CatalogoOutput
}
