package repository

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type FakeAgendamentoRepositorio struct {
	storage map[string]*domain.Agendamento
}

func NovoFakeAgendamentoRepositorio() port.AgendamentoRepositorio {
	return &FakeAgendamentoRepositorio{storage: make(map[string]*domain.Agendamento)}
}

func (r *FakeAgendamentoRepositorio) CriaAgendamento(agendamento *domain.Agendamento) error {
	r.storage[agendamento.ID] = agendamento
	return nil
}
