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

var (
	ErrDuracaoInvalida = errors.New("duração padrão inválida")
	ErrPrecoInvalido   = errors.New("preço inválido")
)

func NovoCatalogo(nome string, duracao int, preco float64, categoria string) (*Catalogo, error) {

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
