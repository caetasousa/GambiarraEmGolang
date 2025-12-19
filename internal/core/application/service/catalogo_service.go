package service

import (
	"errors"
	"fmt"

	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/mapper"
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type CatalogoService struct {
	repo port.CatalogoRepositorio
}

func NovoCatalogoService(r port.CatalogoRepositorio) *CatalogoService {
	return &CatalogoService{repo: r}
}

func (s *CatalogoService) Cadastra(cmd *input.CatalogoInput) (*output.CatalogoOutput, error) {

	catalogo, err := domain.NovoCatalogo(
		cmd.Nome,
		cmd.DuracaoPadrao,
		cmd.Preco,
		cmd.Categoria,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Salvar(catalogo); err != nil {
		return nil, err
	}

	return mapper.FromCatalogo(catalogo), nil
}

var ErrCatalogoNaoEncontrado = errors.New("catálogo não encontrado")
var ErrFalhaInfraestrutura = errors.New("falha na infraestrutura ao buscar catálogo")

func (s *CatalogoService) BuscarPorId(id string) (*output.CatalogoOutput, error) {
	catalogo, err := s.repo.BuscarPorId(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFalhaInfraestrutura, err)
	}
	if catalogo == nil {
		return nil, ErrCatalogoNaoEncontrado
	}
	return mapper.FromCatalogo(catalogo), nil
}
