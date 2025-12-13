package domain

import (
	"errors"

	"github.com/rs/xid"
)

type Catalogo struct {
	ID   string
	Nome string
	// DuracaoPadrao em minutos
	DuracaoPadrao int
	Preco         float64
	Categoria     string
}

func NovoCatalogo(nome string, duracao int, preco float64, categoria string) (*Catalogo, error) {

	if duracao <= 1 {
		return nil, errors.New("duração padrão deve ser maior que 1 minuto")
	}

	if preco < 0 {
		return nil, errors.New("preço não pode ser negativo")
	}

	return &Catalogo{
		ID:            xid.New().String(),
		Nome:          nome,
		DuracaoPadrao: duracao,
		Preco:         preco,
		Categoria:     categoria,
	}, nil
}
