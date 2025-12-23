package domain

import (
	"github.com/rs/xid"
)

type Prestador struct {
	ID       string
	Nome     string
	Cpf      string
	Email    string
	Telefone string
	Ativo    bool
	Catalogo []Catalogo
	Agenda   []AgendaDiaria
}

func NovoPrestador(nome, cpf, email, telefone string, catalogos []Catalogo) (*Prestador, error) {
	if len(catalogos) == 0 {
		return nil, ErrPrestadorDeveTerCatalogo
	}

	return &Prestador{
		ID:       xid.New().String(),
		Nome:     nome,
		Cpf:      cpf,
		Email:    email,
		Telefone: telefone,
		Ativo:    true,
		Catalogo: catalogos,
		Agenda:   []AgendaDiaria{},
	}, nil
}

func (p *Prestador) AdicionarAgenda(agenda *AgendaDiaria) error {
	if !p.Ativo {
		return ErrPrestadorInativo
	}

	for _, a := range p.Agenda {
		if a.Data == agenda.Data {
			return ErrAgendaDuplicada
		}
	}

	p.Agenda = append(p.Agenda, *agenda)
	return nil
}
