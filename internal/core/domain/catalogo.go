package domain

import (
	"github.com/rs/xid"
)

type Catalogo struct {
	ID   string
	Nome string
	// DuracaoPadrao em minutos
	DuracaoPadrao int
	Preco         int
	Categoria     string
}

func NovoCatalogo(nome string, duracao int, preco int, categoria string) (*Catalogo, error) {

	if duracao <= 1 {
		return nil, ErrDuracaoInvalida
	}

	if preco < 0 {
		return nil, ErrPrecoInvalido
	}

	return &Catalogo{
		ID:            xid.New().String(),
		Nome:          nome,
		DuracaoPadrao: duracao,
		Preco:         preco,
		Categoria:     categoria,
	}, nil
}
