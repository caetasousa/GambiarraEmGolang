package service

import (
	"errors"
	"fmt"

	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/mapper"
	"meu-servico-agenda/internal/core/application/output"
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

func (s *PrestadorService) Cadastra(cmd *input.CadastrarPrestadorInput) (*output.CriarPrestadorOutput, error) {

	cpf := cpfcnpj.Clean(cmd.CPF)

	prestadorExistente, err := s.prestadorRepo.BuscarPorCPF(cpf)
	if err != nil {
		return nil, err
	}
	if prestadorExistente != nil {
		return nil, fmt.Errorf("%w: %s", ErrCPFJaCadastrado, cpf)
	}

	catalogos := []domain.Catalogo{}
	for _, id := range cmd.CatalogoIDs {
		c, err := s.catalogoRepo.BuscarPorId(id)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, fmt.Errorf("%w: %s", ErrCatalogoNaoExiste, id)
		}
		catalogos = append(catalogos, *c)
	}

	prestador, err := domain.NovoPrestador(
		cmd.Nome,
		cpf,
		cmd.Email,
		cmd.Telefone,
		catalogos,
	)
	if err != nil {
		return nil, err
	}

	if err := s.prestadorRepo.Salvar(prestador); err != nil {
		return nil, err
	}

	out := mapper.FromDomainToCriarOutput(prestador)

	return out, nil
}

var ErrPrestadorNaoEncontrado = errors.New("prestador não encontrado")

func (s *PrestadorService) AdicionarAgenda(prestadorID string, cmd *input.AdicionarAgendaInput) error {

	prestador, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		return err
	}
	if prestador == nil {
		return ErrPrestadorNaoEncontrado
	}

	intervalos := make([]domain.IntervaloDiario, 0, len(cmd.Intervalos))
	for _, i := range cmd.Intervalos {
		intervalo, err := domain.NovoIntervaloDiario(i.Inicio, i.Fim)
		if err != nil {
			return err
		}
		intervalos = append(intervalos, *intervalo)
	}

	agenda, err := domain.NovaAgendaDiaria(cmd.Data, intervalos)
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

func (s *PrestadorService) BuscarPorId(id string) (*output.BuscarPrestadorOutput, error) {
	prestador, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		return nil, err
	}
	if prestador == nil {
		return nil, ErrPrestadorNaoEncontrado
	}

	out := mapper.FromPrestador(prestador)

	return out, nil
}
