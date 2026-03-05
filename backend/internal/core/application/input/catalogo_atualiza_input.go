package input

type CatalogoUpdateInput struct {
	ID            string
	Nome          string
	DuracaoPadrao int
	Preco         int
	Categoria     string
	ImagemUrl     string
}