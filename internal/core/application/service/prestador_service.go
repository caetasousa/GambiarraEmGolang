package service

import (
	"errors"
	"fmt"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"

	"github.com/klassmann/cpfcnpj"
)

type PrestadorService struct {
	prestadorRepo    port.PrestadorRepositorio
	catalogoRepo     port.CatalogoRepositorio
	agendaDiariaRepo port.AgendaDiariaRepositorio
}

func NovaPrestadorService(pr port.PrestadorRepositorio, cr port.CatalogoRepositorio, ad port.AgendaDiariaRepositorio) *PrestadorService {
	return &PrestadorService{
		prestadorRepo:    pr,
		catalogoRepo:     cr,
		agendaDiariaRepo: ad,
	}
}

var (
	ErrCPFJaCadastrado   = errors.New("cpf já possui um cadastro")
	ErrCatalogoNaoExiste = errors.New("catálogo não existe")
)

func (s *PrestadorService) Cadastra(req *request_prestador.PrestadorRequest) (*domain.Prestador, error) {
	cpf := cpfcnpj.Clean(req.Cpf)
	prestadorExistente, err := s.prestadorRepo.BuscarPorCPF(cpf)
	if err != nil {
		return nil, err
	}
	if prestadorExistente != nil {
		return nil, fmt.Errorf("%w: %s", ErrCPFJaCadastrado, cpf)
	}

	catalogos := []domain.Catalogo{}
	for _, id := range req.CatalogoIDs {
		c, err := s.catalogoRepo.BuscarPorId(id)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, fmt.Errorf("%w: %s", ErrCatalogoNaoExiste, id)
		}
		catalogos = append(catalogos, *c)
	}

	prestador, err := req.ToPrestador(catalogos)
	if err != nil {
		return nil, err
	}

	if err := s.prestadorRepo.Salvar(prestador); err != nil {
		return nil, err
	}

	return prestador, nil
}

var ErrPrestadorNaoEncontrado = errors.New("prestador não encontrado")

func (s *PrestadorService) AdicionarAgenda(prestadorID string, req *request_prestador.AgendaDiariaRequest) error {

	prestador, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		return err
	}

	if prestador == nil {
		return ErrPrestadorNaoEncontrado
	}

	agenda, err := req.ToAgendaDiaria()
	if err != nil {
		return err
	}

	if err := prestador.AdicionarAgenda(agenda); err != nil {
		return err
	}

	if err := s.agendaDiariaRepo.Salvar(agenda); err != nil {
		return err
	}

	return s.prestadorRepo.Salvar(prestador)
}

func (s *PrestadorService) BuscarPorId(id string) (*domain.Prestador, error) {
	prestador, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		return nil, fmt.Errorf("falha na infraestrutura ao buscar prestador: %w", err)
	}
	if prestador == nil {
		return nil, fmt.Errorf("prestador não encontrado")
	}
	return prestador, nil
}
