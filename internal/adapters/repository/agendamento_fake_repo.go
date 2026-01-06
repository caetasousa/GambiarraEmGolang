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

func (r *FakeAgendamentoRepositorio) BuscarAgendamentoPrestadorAPartirDaData(prestadorID string, data time.Time) ([]*domain.Agendamento, error) {
	var resultados []*domain.Agendamento

	for _, agendamento := range r.storage {
		// Verifica se é do prestador correto
		// E se o agendamento inicia na data informada ou depois
		if agendamento.Prestador != nil &&
			agendamento.Prestador.ID == prestadorID &&
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

func (r *FakeAgendamentoRepositorio) ListarPorPrestadorEPeriodoPaginado(prestadorID string, inicio time.Time, fim time.Time, limit, offset int) ([]*domain.Agendamento, error) {
	var resultados []*domain.Agendamento

	for _, agendamento := range r.storage {
		if agendamento.Prestador.ID == prestadorID &&
			inicio.Before(agendamento.DataHoraFim) &&
			fim.After(agendamento.DataHoraInicio) {
			resultados = append(resultados, agendamento)
		}
	}

	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].DataHoraInicio.Before(resultados[j].DataHoraInicio)
	})

	start := offset
	if start > len(resultados) {
		return []*domain.Agendamento{}, nil
	}

	end := offset + limit
	if end > len(resultados) {
		end = len(resultados)
	}

	return resultados[start:end], nil
}

func (r *FakeAgendamentoRepositorio) ListarPorClienteEPeriodoPaginado(clienteID string, inicio time.Time, fim time.Time, limit, offset int) ([]*domain.Agendamento, error) {
	var resultados []*domain.Agendamento

	for _, agendamento := range r.storage {
		if agendamento.Cliente.ID == clienteID &&
			inicio.Before(agendamento.DataHoraFim) &&
			fim.After(agendamento.DataHoraInicio) {
			resultados = append(resultados, agendamento)
		}
	}

	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].DataHoraInicio.Before(resultados[j].DataHoraInicio)
	})

	start := offset
	if start > len(resultados) {
		return []*domain.Agendamento{}, nil
	}

	end := offset + limit
	if end > len(resultados) {
		end = len(resultados)
	}

	return resultados[start:end], nil
}

func (r *FakeAgendamentoRepositorio) ContarPorPrestadorEPeriodo(prestadorID string, inicio time.Time, fim time.Time) (int, error) {
	count := 0

	for _, agendamento := range r.storage {
		if agendamento.Prestador.ID == prestadorID &&
			inicio.Before(agendamento.DataHoraFim) &&
			fim.After(agendamento.DataHoraInicio) {
			count++
		}
	}

	return count, nil
}

func (r *FakeAgendamentoRepositorio) ContarPorClienteEPeriodo(clienteID string, inicio time.Time, fim time.Time) (int, error) {
	count := 0

	for _, agendamento := range r.storage {
		if agendamento.Cliente.ID == clienteID &&
			inicio.Before(agendamento.DataHoraFim) &&
			fim.After(agendamento.DataHoraInicio) {
			count++
		}
	}

	return count, nil
}