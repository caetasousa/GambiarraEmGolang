package domain

import "time"

type CatalogoServico struct {
	ID            string 
	Nome          string
	DuracaoPadrao time.Duration 
	Preco         float64
	Categoria     string
}
