package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (s *PrestadorService) ListarPrestadores(input *input.PrestadorListInput) ([]*output.BuscarPrestadorOutput, int, error) {
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

	// Buscar prestadores (sempre com filtro)
	prestadores, err := s.prestadorRepo.Listar(input)
	if err != nil {
		return nil, 0, err
	}

	// Contar total (sempre com filtro)
	total, err := s.prestadorRepo.Contar(input.Ativo)
	if err != nil {
		return nil, 0, err
	}

	outputs := mapper.PrestadoresFromDomainOutput(prestadores)

	return outputs, total, nil
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

	err = s.prestadorRepo.AtualizarStatus(id, true)
	if err != nil {
		return err
	}
	return nil
}

func (s *PrestadorService) SalvarAgenda(cmd *input.AdicionarAgendaInput) error {
	// 1. Buscar prestador
	prestador, err := s.prestadorRepo.BuscarPorId(cmd.PrestadorID)
	if err != nil {
		return ErrPrestadorNaoEncontrado
	}

	// 2. Validar se prestador está ativo
	if !prestador.Ativo {
		return ErrPrestadorInativo
	}

	// 3. Criar intervalos usando construtor do domínio
	intervalos := make([]domain.IntervaloDiario, 0, len(cmd.Intervalos))
	for _, i := range cmd.Intervalos {
		intervalo, err := domain.NovoIntervaloDiario(i.Inicio, i.Fim)
		if err != nil {
			return err
		}
		intervalos = append(intervalos, *intervalo)
	}

	// 4. Criar agenda usando construtor do domínio
	novaAgenda, err := domain.NovaAgendaDiaria(cmd.Data, intervalos)
	if err != nil {
		return err
	}

	// 5. Verificar se já existe agenda para essa data
	dataFormatada := cmd.Data.Format("2006-01-02")
	agendaExistente, err := s.agendaDiariaRepo.BuscarAgendaDoDia(cmd.PrestadorID, dataFormatada)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// 6. Se existe, ATUALIZAR. Se não existe, CRIAR
	if agendaExistente != nil {
		// ATUALIZAÇÃO
		novaAgenda.Id = agendaExistente.Id
		return s.agendaDiariaRepo.AtualizarAgenda(novaAgenda, cmd.PrestadorID)
	} else {
		// CRIAÇÃO
		if err := prestador.AdicionarAgenda(novaAgenda); err != nil {
			return err
		}
		return s.agendaDiariaRepo.Salvar(novaAgenda, cmd.PrestadorID)
	}
}

func (s *PrestadorService) DeletarAgenda(prestadorID string, data string) error {
	// 1. Buscar prestador
	prestador, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		return ErrPrestadorNaoEncontrado
	}

	// 2. Validar se prestador está ativo
	if !prestador.Ativo {
		return ErrPrestadorInativo
	}

	// 3. Verificar se agenda existe
	_, err = s.agendaDiariaRepo.BuscarAgendaDoDia(prestadorID, data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAgendaNaoEncontrada
		}
		return err
	}

	// 4. TODO: Verificar se não tem agendamentos para esta data
	// (implementar quando tiver a entidade Agendamento)

	// 5. Remover agenda do domínio
	if err := prestador.RemoverAgenda(data); err != nil {
		return err
	}

	// 6. Deletar do repositório
	return s.agendaDiariaRepo.DeletarAgenda(prestadorID, data)
}

func (s *PrestadorService) BuscarPrestadoresDisponiveisPorData(input *input.PrestadorListDataInput) ([]*output.BuscarPrestadorOutput, int, error) {

	// Validar formato da data
	dataTime, err := time.Parse("2006-01-02", input.Data)
	if err != nil {
		return nil, 0, ErrFormatoDataInvalido
	}

	// Validar se a data não está no passado
	if err := domain.ValidarDataNoPassado(dataTime); err != nil {
		return nil, 0, err
	}

	// Buscar prestadores disponíveis
	prestadores, err := s.prestadorRepo.BuscarPrestadoresDisponiveisPorData(input.Data, input.Page, input.Limit)
	if err != nil {
		return nil, 0, ErrAoBuscarPrestadoresDisponiveis
	}

	// Contar total
	total, err := s.prestadorRepo.ContarPrestadoresDisponiveisPorData(input.Data)
	if err != nil {
		return nil, 0, ErrAoContarPrestadoresDisponiveis
	}

	// Mapear para output
	outputs := mapper.PrestadoresFromDomainOutput(prestadores)

	return outputs, total, nil
}
