package domain

type Catalogo struct {
	ID            string
	Nome          string
	// DuracaoPadrao em minutos
	DuracaoPadrao int
	Preco         float64
	Categoria     string
}
