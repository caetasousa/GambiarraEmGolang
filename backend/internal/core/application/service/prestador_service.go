package service

import (
	"errors"

	"time"

	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/mapper"
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
)

type PrestadorService struct {
	prestadorRepo    port.PrestadorRepositorio
	catalogoRepo     port.CatalogoRepositorio
	agendaDiariaRepo port.AgendaDiariaRepositorio
	authService      *AuthService
	clienteRepo      port.ClienteRepositorio
}

func NovaPrestadorService(pr port.PrestadorRepositorio, cr port.CatalogoRepositorio, ad port.AgendaDiariaRepositorio, authService *AuthService, clienteRepo port.ClienteRepositorio) *PrestadorService {
	return &PrestadorService{
		prestadorRepo:    pr,
		catalogoRepo:     cr,
		agendaDiariaRepo: ad,
		authService:      authService,
		clienteRepo:      clienteRepo,
	}
}

func (s *PrestadorService) Cadastra(cmd *input.CadastrarPrestadorInput) (*output.CriarPrestadorOutput, error) {
	// CPF já vem validado e limpo do Request Mapper via Value Object (mas aqui usamos a string do input)
	cpf := cmd.CPF

	// BuscarPorCPF retorna (nil, nil) se não encontrar - isso é esperado para novo cadastro
	prestadorExistente, err := s.prestadorRepo.BuscarPorCPF(cpf)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	if prestadorExistente != nil {
		return nil, ErrCPFJaCadastrado
	}

	// Validar se o email já está sendo usado por um cliente
	if s.clienteRepo != nil {
		clienteExistente, err := s.clienteRepo.BuscarPorEmail(cmd.Email)
		if err != nil {
			return nil, ErrFalhaInfraestrutura
		}
		if clienteExistente != nil {
			return nil, ErrEmailJaUsadoPorCliente
		}
	}

	catalogos := []domain.Catalogo{}
	for _, id := range cmd.CatalogoIDs {
		c, err := s.catalogoRepo.BuscarPorId(id)
		if err != nil {
			if errors.Is(err, repository.ErrCatalogoNaoEncontrado) {
				return nil, ErrCatalogoNaoEncontrado
			}
			return nil, ErrFalhaInfraestrutura
		}
		catalogos = append(catalogos, *c)
	}

	// Criar prestador (senha já vem hasheada do request)
	prestador, err := domain.NovoPrestador(
		cmd.ID,
		cmd.Nome,
		cpf,
		cmd.Email,
		cmd.Telefone,
		cmd.Senha, // Senha já é o hash
		cmd.ImagemUrl,
		catalogos,
	)
	if err != nil {
		return nil, err
	}

	if err := s.prestadorRepo.Salvar(prestador); err != nil {
		if errors.Is(err, repository.ErrCPFDuplicado) {
			return nil, ErrCPFJaCadastrado
		}
		return nil, ErrFalhaInfraestrutura
	}

	out := mapper.FromDomainToCriarOutput(prestador)
	return out, nil
}

func (s *PrestadorService) BuscarPorId(id string) (*output.BuscarPrestadorOutput, error) {
	prestador, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return nil, ErrPrestadorNaoEncontrado
		}
		return nil, ErrFalhaInfraestrutura
	}

	out := mapper.FromPrestador(prestador)
	return out, nil
}

func (s *PrestadorService) Atualizar(input *input.AlterarPrestadorInput) error {
	if len(input.CatalogoIDs) == 0 {
		return domain.ErrPrestadorDeveTerCatalogo
	}

	if err := s.prestadorRepo.Atualizar(input); err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		if errors.Is(err, repository.ErrCatalogoNaoEncontrado) {
			return ErrCatalogoNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	return nil
}

func (s *PrestadorService) ListarPrestadores(input *input.PrestadorListInput) ([]*output.BuscarPrestadorOutput, int, error) {
	prestadores, err := s.prestadorRepo.Listar(input)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}

	total, err := s.prestadorRepo.Contar(input.Ativo)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}

	outputs := mapper.PrestadoresFromDomainOutput(prestadores)
	return outputs, total, nil
}

func (s *PrestadorService) Inativar(id string) error {
	_, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	if err := s.prestadorRepo.AtualizarStatus(id, false); err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	return nil
}

func (s *PrestadorService) Ativar(id string) error {
	_, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	if err := s.prestadorRepo.AtualizarStatus(id, true); err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	return nil
}

func (s *PrestadorService) SalvarAgenda(cmd *input.AdicionarAgendaInput) error {
	prestador, err := s.prestadorRepo.BuscarPorId(cmd.PrestadorID)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	if !prestador.Ativo {
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

	novaAgenda, err := domain.NovaAgendaDiaria(cmd.Data, intervalos)
	if err != nil {
		return err
	}

	// BuscarAgendaDoDia retorna (nil, nil) se não encontrar - isso é esperado
	dataFormatada := cmd.Data.Format("2006-01-02")
	agendaExistente, err := s.agendaDiariaRepo.BuscarAgendaDoDia(cmd.PrestadorID, dataFormatada)
	if err != nil {
		return ErrFalhaInfraestrutura
	}

	if agendaExistente != nil {
		// ATUALIZAÇÃO
		novaAgenda.Id = agendaExistente.Id
		if err := s.agendaDiariaRepo.AtualizarAgenda(novaAgenda, cmd.PrestadorID); err != nil {
			return ErrFalhaInfraestrutura
		}
		return nil
	}

	// CRIAÇÃO
	if err := prestador.AdicionarAgenda(novaAgenda); err != nil {
		return err
	}

	if err := s.agendaDiariaRepo.Salvar(novaAgenda, cmd.PrestadorID); err != nil {
		return ErrFalhaInfraestrutura
	}

	return nil
}

func (s *PrestadorService) DeletarAgenda(prestadorID string, data string) error {
	prestador, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return ErrPrestadorNaoEncontrado
		}
		return ErrFalhaInfraestrutura
	}

	if !prestador.Ativo {
		return ErrPrestadorInativo
	}

	// Verificar se agenda existe
	agendaExistente, err := s.agendaDiariaRepo.BuscarAgendaDoDia(prestadorID, data)
	if err != nil {
		return ErrFalhaInfraestrutura
	}
	if agendaExistente == nil {
		return ErrAgendaNaoEncontrada
	}

	if err := prestador.RemoverAgenda(data); err != nil {
		return err
	}

	if err := s.agendaDiariaRepo.DeletarAgenda(prestadorID, data); err != nil {
		if errors.Is(err, repository.ErrAgendaNaoEncontrada) {
			return ErrAgendaNaoEncontrada
		}
		return ErrFalhaInfraestrutura
	}

	return nil
}

func (s *PrestadorService) BuscarPrestadoresDisponiveisPorData(input *input.PrestadorListDataInput) ([]*output.BuscarPrestadorOutput, int, error) {
	dataTime, err := time.Parse("2006-01-02", input.Data)
	if err != nil {
		return nil, 0, ErrFormatoDataInvalido
	}

	if err := domain.ValidarDataNoPassado(dataTime); err != nil {
		return nil, 0, err
	}

	prestadores, err := s.prestadorRepo.BuscarPrestadoresDisponiveisPorData(input.Data, input.Limit, input.Offset)
	if err != nil {
		return nil, 0, ErrAoBuscarPrestadoresDisponiveis
	}

	total, err := s.prestadorRepo.ContarPrestadoresDisponiveisPorData(input.Data)
	if err != nil {
		return nil, 0, ErrAoContarPrestadoresDisponiveis
	}

	outputs := mapper.PrestadoresFromDomainOutput(prestadores)
	return outputs, total, nil
}
