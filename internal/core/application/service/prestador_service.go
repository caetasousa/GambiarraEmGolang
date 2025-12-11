package service

import (
	"meu-servico-agenda/internal/adapters/http/prestador/request"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorService struct {
	repo port.PrestadorDeServicoRepositorio
}

func NovoPrestadorService(r port.PrestadorDeServicoRepositorio) *PrestadorService {
	return &PrestadorService{repo: r}
}

func (s *PrestadorService) Cadastra(input request.PrestadorRequest) (*domain.Prestador, error) {
	prestador := input.ToPrestador()
	err := s.repo.Salvar(prestador)
	if err != nil {
		return nil, err
	}
	return prestador, nil
}
