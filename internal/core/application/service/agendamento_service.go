package service

import (
	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type AgendamentoService struct {
	prestadorRepo   port.PrestadorRepositorio
	agendamentoRepo port.AgendamentoRepositorio
}

func NovaAgendamentoService(pr port.PrestadorRepositorio, ar port.AgendamentoRepositorio) *AgendamentoService {
	return &AgendamentoService{
		prestadorRepo:   pr,
		agendamentoRepo: ar,
	}
}

func (s *AgendamentoService) CadastraAgendamento(input request_agendamento.AgendamentoRequest) (*domain.Agendamento, error) {

	agendamento, err := input.ToAgendamento()
	if err != nil {
		return nil, err
	}

	if err := s.agendamentoRepo.CriaAgendamento(agendamento); err != nil {
		return nil, err
	}

	return agendamento, nil
}
