package repository

import (
	"meu-servico-agenda/internal/core/domain"
)

type AgendaFakeRepo struct {
	Agendas map[string]*domain.AgendaDiaria
}

func NovoAgendaFakeRepo() *AgendaFakeRepo {
	return &AgendaFakeRepo{
		Agendas: make(map[string]*domain.AgendaDiaria),
	}
}

func (r *AgendaFakeRepo) Salvar(agenda *domain.AgendaDiaria) error {
	r.Agendas[agenda.Id] = agenda
	return nil
}

func (r *AgendaFakeRepo) BuscarPorId(id string) (*domain.AgendaDiaria, error) {
	agenda, ok := r.Agendas[id]
	if !ok {
		return nil, nil
	}
	return agenda, nil
}
