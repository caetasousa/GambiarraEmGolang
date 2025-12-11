package service

import (
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorService struct {
	repo port.PrestadorDeServicoRepositorio
}

func NovoPrestadorService(r port.PrestadorDeServicoRepositorio) *PrestadorService {
	return &PrestadorService{repo: r}
}

func (s *PrestadorService) Cadastra(prestador *domain.Prestador) (*domain.Prestador, error) {
	if err := s.repo.Salvar(prestador); err != nil {
		return nil, err
	}
	return prestador, nil
}
