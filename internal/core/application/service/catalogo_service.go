package service

import (
	"errors"
	"fmt"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoService struct {
	repo port.CatalogoRepositorio
}

func NovoCatalogoService(r port.CatalogoRepositorio) *CatalogoService {
	return &CatalogoService{repo: r}
}

func (s *CatalogoService) Cadastra(input *domain.Catalogo) (*domain.Catalogo, error) {
	if err := s.repo.Salvar(input); err != nil {
		return nil, err
	}

	return input, nil
}

func (s *CatalogoService) BuscarPorId(id string) (*domain.Catalogo, error) {
	catalogo, err := s.repo.BuscarPorId(id)
	if err != nil {
		return nil, fmt.Errorf("falha na infraestrutura ao buscar catalogo: %w", err)
	}
	if catalogo == nil {
		return nil, errors.New("catalogo nao encontrado") // <- sem wrap
	}
	return catalogo, nil
}
