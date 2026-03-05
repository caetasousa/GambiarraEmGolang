package service

import (
	"errors"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/mapper"
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"time"
)

type AgendamentoService struct {
	prestadorRepo   port.PrestadorRepositorio
	agendamentoRepo port.AgendamentoRepositorio
	catalogoRepo    port.CatalogoRepositorio
	clienteRepo     port.ClienteRepositorio
}

func NovaAgendamentoService(pr port.PrestadorRepositorio, ar port.AgendamentoRepositorio, cr port.CatalogoRepositorio, cl port.ClienteRepositorio) *AgendamentoService {
	return &AgendamentoService{
		prestadorRepo:   pr,
		agendamentoRepo: ar,
		catalogoRepo:    cr,
		clienteRepo:     cl,
	}
}

func (s *AgendamentoService) CadastraAgendamento(input *input.CadastrarAgendamentoInput) (*output.AgendamentoOutput, error) {
	if err := domain.ValidarDataNoPassado(input.DataHoraInicio); err != nil {
		return nil, domain.ErrDataEstaNoPassado
	}

	// Buscar cliente
	cliente, err := s.clienteRepo.BuscarPorId(input.ClienteID)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}
	if cliente == nil {
		return nil, ErrClienteNaoEncontrado
	}

	// Buscar prestador
	prestador, err := s.prestadorRepo.BuscarPorId(input.PrestadorID)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return nil, ErrPrestadorNaoEncontrado
		}
		return nil, ErrFalhaInfraestrutura
	}

	// Buscar catálogo
	catalogo, err := s.catalogoRepo.BuscarPorId(input.CatalogoID)
	if err != nil {
		if errors.Is(err, repository.ErrCatalogoNaoEncontrado) {
			return nil, ErrCatalogoNaoEncontrado
		}
		return nil, ErrFalhaInfraestrutura
	}

	dataHorarioFim := input.DataHoraInicio.Add(time.Duration(catalogo.DuracaoPadrao) * time.Minute)

	// Extrai apenas a data (sem hora) para validações por dia
	inicioDoDia := time.Date(
		input.DataHoraInicio.Year(),
		input.DataHoraInicio.Month(),
		input.DataHoraInicio.Day(),
		0, 0, 0, 0,
		input.DataHoraInicio.Location(),
	)
	fimDoDia := inicioDoDia.Add(24 * time.Hour)

	// Valida se já existe agendamento da mesma categoria no mesmo dia
	agendamentosDoDia, err := s.agendamentoRepo.BuscarPorClienteEPeriodo(
		input.ClienteID,
		inicioDoDia,
		fimDoDia,
	)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	for _, agend := range agendamentosDoDia {
		if agend.Catalogo.ID == catalogo.ID {
			return nil, ErrAgendamentoDuplo
		}
	}

	// Busca a agenda do prestador para o dia solicitado
	dia := input.DataHoraInicio.Format("2006-01-02")
	agendaDoDia, err := s.prestadorRepo.BuscarAgendaDoDia(prestador.ID, dia)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	// Valida se o dia é atendido pelo prestador
	if agendaDoDia == nil {
		return nil, ErrDiaIndisponivel
	}

	// Valida se o horário solicitado está dentro dos horários disponíveis do dia
	if !agendaDoDia.PermiteAgendamento(input.DataHoraInicio, dataHorarioFim) {
		return nil, ErrHorarioIndisponivel
	}

	// Um prestador não pode ter dois atendimentos no mesmo período
	conflitosPrestador, err := s.agendamentoRepo.BuscarPorPrestadorEPeriodo(prestador.ID, input.DataHoraInicio, dataHorarioFim)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}
	if len(conflitosPrestador) > 0 {
		return nil, ErrPrestadorOcupado
	}

	// Um cliente não pode ter dois agendamentos simultâneos
	conflitosCliente, err := s.agendamentoRepo.BuscarPorClienteEPeriodo(cliente.ID, input.DataHoraInicio, dataHorarioFim)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}
	if len(conflitosCliente) > 0 {
		return nil, ErrClienteOcupado
	}

	agendamento, err := domain.NovoAgendamento(
		input.ID,
		cliente,
		prestador,
		catalogo,
		input.DataHoraInicio,
		dataHorarioFim,
		input.Notas,
	)
	if err != nil {
		return nil, err
	}

	if err := s.agendamentoRepo.CriaAgendamento(agendamento); err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	out := mapper.NovoAgendamentoOutput(agendamento)
	return out, nil
}

func (s *AgendamentoService) ConsultaAgendamentoClienteData(request input.AgendamentoDataInput, id string) ([]*output.AgendamentoOutput, error) {
	cliente, err := s.clienteRepo.BuscarPorId(id)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}
	if cliente == nil {
		return nil, ErrClienteNaoEncontrado
	}

	agendamentos, err := s.agendamentoRepo.BuscarAgendamentoClienteAPartirDaData(id, request.Data)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	out := mapper.BuscaAgendamentoData(agendamentos)
	return out, nil
}

func (s *AgendamentoService) ConsultaAgendamentoPrestadorData(request input.AgendamentoDataInput, id string) ([]*output.AgendamentoOutput, error) {
	_, err := s.prestadorRepo.BuscarPorId(id)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return nil, ErrPrestadorNaoEncontrado
		}
		return nil, ErrFalhaInfraestrutura
	}

	agendamentos, err := s.agendamentoRepo.BuscarAgendamentoPrestadorAPartirDaData(id, request.Data)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	out := mapper.BuscaAgendamentoData(agendamentos)
	return out, nil
}

func (s *AgendamentoService) ListarAgendamentosPrestadorPorPeriodo(prestadorID string, in *input.ListarAgendamentoPrestadorPeriodoInput) ([]*output.AgendamentoOutput, int, error) {
	_, err := s.prestadorRepo.BuscarPorId(prestadorID)
	if err != nil {
		if errors.Is(err, repository.ErrPrestadorNaoEncontrado) {
			return nil, 0, ErrPrestadorNaoEncontrado
		}
		return nil, 0, ErrFalhaInfraestrutura
	}

	agendamentos, err := s.agendamentoRepo.ListarPorPrestadorEPeriodoPaginado(prestadorID, in.DataInicio, in.DataFim, in.Limit, in.Offset)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}

	total, err := s.agendamentoRepo.ContarPorPrestadorEPeriodo(prestadorID, in.DataInicio, in.DataFim)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}

	out := mapper.BuscaAgendamentoData(agendamentos)
	return out, total, nil
}

func (s *AgendamentoService) ListarAgendamentosClientePorPeriodo(clienteID string, in *input.ListarAgendamentoClientePeriodoInput) ([]*output.AgendamentoOutput, int, error) {
	cliente, err := s.clienteRepo.BuscarPorId(clienteID)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}
	if cliente == nil {
		return nil, 0, ErrClienteNaoEncontrado
	}

	agendamentos, err := s.agendamentoRepo.ListarPorClienteEPeriodoPaginado(clienteID, in.DataInicio, in.DataFim, in.Limit, in.Offset)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}

	total, err := s.agendamentoRepo.ContarPorClienteEPeriodo(clienteID, in.DataInicio, in.DataFim)
	if err != nil {
		return nil, 0, ErrFalhaInfraestrutura
	}

	out := mapper.BuscaAgendamentoData(agendamentos)
	return out, total, nil
}
