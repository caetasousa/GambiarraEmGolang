package service

import (
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

func (s *CatalogoService) Cadastra(input *input.CatalogoInput) (*output.CatalogoOutput, error) {

	catalogo, err := domain.NovoCatalogo(
		input.Nome,
		input.DuracaoPadrao,
		input.Preco,
		input.Categoria,
		input.ImagemUrl,
	)

	if err != nil {
		return nil, err
	}

	if err := s.repo.Salvar(catalogo); err != nil {
		return nil, err
	}

	return mapper.FromCatalogoOutput(catalogo), nil
}

func (s *CatalogoService) BuscarPorId(id string) (*output.CatalogoOutput, error) {
	catalogo, err := s.repo.BuscarPorId(id)
	if err != nil {
		return nil, ErrCatalogoNaoEncontrado
	}
	return mapper.FromCatalogoOutput(catalogo), nil
}
