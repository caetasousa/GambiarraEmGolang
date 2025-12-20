package repository

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type FakeAgendaDiariaRepositorio struct {
	storage map[string]*domain.AgendaDiaria // agendaID → agenda
}

func NovoFakeAgendaDiariaRepositorio() port.AgendaDiariaRepositorio {
	return &FakeAgendaDiariaRepositorio{
		storage: make(map[string]*domain.AgendaDiaria),
	}
}

func (r *FakeAgendaDiariaRepositorio) Salvar(agenda *domain.AgendaDiaria) error {
	r.storage[agenda.Id] = agenda
	return nil
}

