package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (s *PrestadorService) Cadastra(cmd *input.CadastrarPrestadorInput) (*output.CriarPrestadorOutput, error) {

	cpf := cpfcnpj.Clean(cmd.CPF)

	prestadorExistente, err := s.prestadorRepo.BuscarPorCPF(cpf)
	if err != nil {
		return nil, err
	}
	//log.Printf("✅ prestador de cpf %s", cmd.CPF)
	if prestadorExistente != nil {
		return nil, fmt.Errorf("%w: %s", ErrCPFJaCadastrado, cpf)
	}

	catalogos := []domain.Catalogo{}
	for _, id := range cmd.CatalogoIDs {
		c, err := s.catalogoRepo.BuscarPorId(id)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrCatalogoNaoExiste, id)
		}
		catalogos = append(catalogos, *c)
	}

	prestador, err := domain.NovoPrestador(
		cmd.Nome,
		cpf,
		cmd.Email,
		cmd.Telefone,
		cmd.ImagemUrl,
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

func (s *PrestadorService) AdicionarAgenda(prestadorID string, cmd *input.AdicionarAgendaInput) error {

	prestador, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		return ErrPrestadorNaoEncontrado
	}

	if prestador.Ativo == false {
		return ErrPrestadorInativo
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

	if err := s.agendaDiariaRepo.Salvar(agenda, prestadorID); err != nil {
		return err
	}
	return nil
}

func (s *PrestadorService) BuscarPorId(id string) (*output.BuscarPrestadorOutput, error) {
	prestador, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		return nil, ErrPrestadorNaoEncontrado
	}

	out := mapper.FromPrestador(prestador)

	return out, nil
}

func (s *PrestadorService) Atualizar(input *input.AlterarPrestadorInput) error {
	if len(input.CatalogoIDs) == 0 {
		return domain.ErrPrestadorDeveTerCatalogo
	}

	if err := s.prestadorRepo.Atualizar(input); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPrestadorNaoEncontrado
		}
		if strings.Contains(err.Error(), "catálogo") && strings.Contains(err.Error(), "não existe") {
			return ErrCatalogoNaoExiste
		}

		return err
	}

	return nil
}

func (s *PrestadorService) Listar(input *input.PrestadorListInput) ([]*output.BuscarPrestadorOutput, int, error) {
	// Validações
	if input.Page <= 0 {
		input.Page = 1
	}

	if input.Limit <= 0 {
		input.Limit = 10
	}

	if input.Limit > 100 {
		input.Limit = 100
	}

	// Busca prestadores do repositório
	prestadores, err := s.prestadorRepo.Listar(input)

	if err != nil {
		return nil, 0, err
	}

	// Busca o total de prestadores
	total, err := s.prestadorRepo.Contar()
	if err != nil {
		return nil, 0, err
	}

	// Converte de domain para output usando mapper
	prestadoresOutput := mapper.PrestadoresFromDomainOutput(prestadores)

	return prestadoresOutput, total, nil
}


// Service
func (s *PrestadorService) Inativar(id string) error {
	prestador, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPrestadorNaoEncontrado
		}
		return err
	}

	prestador.Ativo = false
	
	// Usar método específico no repo ou o Salvar existente
	return s.prestadorRepo.AtualizarStatus(id, false)
}

func (s *PrestadorService) Ativar(id string) error {
	prestador, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPrestadorNaoEncontrado
		}
		return err
	}

	prestador.Ativo = true
	
	return s.prestadorRepo.AtualizarStatus(id, true)
}