package domain

type Catalogo struct {
	ID   string
	Nome string
	// DuracaoPadrao em minutos
	DuracaoPadrao int
	Preco         int
	ImagemUrl     string
	Categoria     string
}

func NovoCatalogo(id, nome string, duracao int, preco int, categoria string, image_url string) (*Catalogo, error) {

	if duracao <= 1 {
		return nil, ErrDuracaoInvalida
	}

	if preco < 0 {
		return nil, ErrPrecoInvalido
	}

	return &Catalogo{
		ID:            id,
		Nome:          nome,
		DuracaoPadrao: duracao,
		Preco:         preco,
		Categoria:     categoria,
		ImagemUrl:     image_url,
	}, nil
}
