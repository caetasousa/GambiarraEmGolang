package domain

import (
	"errors"
	"time"

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
		return nil, errors.New("prestador deve ter ao menos um catálogo de serviços")
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

var (
	ErrAgendaDuplicada  = errors.New("agenda duplicada")
	ErrPrestadorInativo = errors.New("prestador inativo")
)

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

func (p *Prestador) PodeAgendar(inicio, fim time.Time) bool {
	data := inicio.Format("2006-01-02")

	for _, agenda := range p.Agenda {
		if agenda.Data == data {
			return agenda.PermiteAgendamento(inicio, fim)
		}
	}

	return false
}
