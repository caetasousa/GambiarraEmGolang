package repository

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"sort"
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
func (r *FakeAgendamentoRepositorio) BuscarAgendamentoClienteAPartirDaData(clienteID string, data time.Time) ([]*domain.Agendamento, error) {
	var resultados []*domain.Agendamento

	for _, agendamento := range r.storage {
		// Verifica se é do cliente correto
		// E se o agendamento inicia na data informada ou depois
		if agendamento.Cliente != nil && 
		   agendamento.Cliente.ID == clienteID &&
		   !agendamento.DataHoraInicio.Before(data) {
			resultados = append(resultados, agendamento)
		}
	}

	// Ordena por data/hora de início (do mais antigo para o mais recente)
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].DataHoraInicio.Before(resultados[j].DataHoraInicio)
	})

	return resultados, nil
}

func (r *FakeAgendamentoRepositorio) BuscarAgendamentoPrestadorAPartirDaData(clienteID string, data time.Time) ([]*domain.Agendamento, error) {
	return nil,nil
}