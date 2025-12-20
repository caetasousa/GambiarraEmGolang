package repository

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"time"
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

func (r *FakeAgendamentoRepositorio) BuscarPorPrestadorEPeriodo(prestadorID string, inicio, fim time.Time) ([]*domain.Agendamento, error) {

	var resultados []*domain.Agendamento

	for _, agendamento := range r.storage {
		if agendamento.Prestador.ID == prestadorID &&
			inicio.Before(agendamento.DataHoraFim) &&
			fim.After(agendamento.DataHoraInicio) {

			resultados = append(resultados, agendamento)
		}
	}

	return resultados, nil
}

func (r *FakeAgendamentoRepositorio) BuscarPorClienteEPeriodo(clienteID string, inicio, fim time.Time) ([]*domain.Agendamento, error) {

	var resultados []*domain.Agendamento

	for _, agendamento := range r.storage {
		if agendamento.Cliente.ID == clienteID &&
			inicio.Before(agendamento.DataHoraFim) &&
			fim.After(agendamento.DataHoraInicio) {

			resultados = append(resultados, agendamento)
		}
	}

	return resultados, nil
}
