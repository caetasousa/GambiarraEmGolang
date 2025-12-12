package service

import (
	"fmt"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type ServiceAgendaDiaria struct {
	repo port.AgendaDiariaRepositorio
}

func NovaServiceAgendaDiaria(r port.AgendaDiariaRepositorio) *ServiceAgendaDiaria {
	return &ServiceAgendaDiaria{repo: r}
}

func (s *ServiceAgendaDiaria) Cadastra(input *domain.AgendaDiaria) (*domain.AgendaDiaria, error) {

	if err := s.repo.Salvar(input); err != nil {
		return nil, fmt.Errorf("falha ao salvar agenda: %w", err)
	}

	return input, nil
}
