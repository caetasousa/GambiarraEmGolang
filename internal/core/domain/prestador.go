package domain

import (
	"errors"

	"github.com/rs/xid"
)

type Prestador struct {
	ID       string
	Nome     string
	Email    string
	Telefone string
	Ativo    bool
	Catalogo []Catalogo
	Agenda   []AgendaDiaria
}

func NovoPrestador(nome, email, telefone string, catalogos []Catalogo) (*Prestador, error) {
	if len(catalogos) == 0 {
		return nil, errors.New("prestador deve ter ao menos um catálogo de serviços")
	}

	return &Prestador{
		ID:       xid.New().String(),
		Nome:     nome,
		Email:    email,
		Telefone: telefone,
		Ativo:    true,
		Catalogo: catalogos,
		Agenda:   []AgendaDiaria{},
	}, nil
}
