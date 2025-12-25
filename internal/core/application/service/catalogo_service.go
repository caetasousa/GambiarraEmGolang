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

func (s *CatalogoService) Listar(in *input.ListCatalogoInput) ([]*output.CatalogoOutput, int, error) {

	page := in.Page
	limit := in.Limit

	if page <= 0 {
		page = 1
	}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	catalogos, err := s.repo.Listar(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Contar()
	if err != nil {
		return nil, 0, err
	}

	return mapper.CatalogosFromDomainOutput(catalogos), total, nil
}

func (s *CatalogoService) Atualizar(input *input.CatalogoUpdateInput) error {
	// Verifica se existe
	catalogo, err := s.repo.BuscarPorId(input.ID)
	if err != nil {
		return ErrCatalogoNaoEncontrado
	}

	// objeto para update
	catalogo.Nome = input.Nome
	catalogo.Categoria = input.Categoria

	if input.DuracaoPadrao != 0 && input.DuracaoPadrao <= 1 {
		return domain.ErrDuracaoInvalida
	}
	if input.DuracaoPadrao > 1 {
		catalogo.DuracaoPadrao = input.DuracaoPadrao
	}

	if input.Preco < 0 {
		return domain.ErrPrecoInvalido
	}
	if input.Preco > 0 {
		catalogo.Preco = input.Preco
	}

	catalogo.ImagemUrl = input.ImagemUrl

	if err := s.repo.Atualizar(catalogo); err != nil {
		return err
	}

	return nil
}

func (s *CatalogoService) Deletar(id string) error {
	// Verifica se o catálogo existe
	_, err := s.repo.BuscarPorId(id)
	if err != nil {
		return ErrCatalogoNaoEncontrado
	}

	// Deleta o catálogo
	if err := s.repo.Deletar(id); err != nil {
		return ErrFalhaInfraestrutura
	}

	return nil
}