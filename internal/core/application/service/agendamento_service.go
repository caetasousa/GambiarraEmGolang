package service

import (
	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
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

func (s *AgendamentoService) CadastraAgendamento(request request_agendamento.AgendamentoRequest) (*output.AgendamentoOutput, error) {
	input, err := request.ToAgendamento()

	cliente, err := s.clienteRepo.BuscarPorId(input.ClienteID)
	if err != nil || cliente == nil {
		return nil, ErrClienteNaoExiste
	}

	prestador, err := s.prestadorRepo.BuscarPorId(input.PrestadorID)
	if err != nil || prestador == nil {
		return nil, ErrPrestadorNaoExiste
	}

	catalogo, err := s.catalogoRepo.BuscarPorId(input.CatalogoID)
	if err != nil || catalogo == nil {
		return nil, ErrCatalogoNaoExiste
	}

	dataHorarioFim := input.DataHoraInicio.Add(time.Duration(catalogo.DuracaoPadrao) * time.Minute)

	// Busca a agenda do prestador para o dia solicitado
	// Aqui validamos se o prestador trabalha nesse dia específico
	agendaDoDia, err := s.prestadorRepo.BuscarAgendaDoDia(prestador.ID, input.DataHoraInicio.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	// Valida se o dia é atendido pelo prestador
	// Exemplo: prestador não trabalha à tarde ou não trabalha nesse dia
	if agendaDoDia == nil {
		return nil, ErrDiaIndisponivel
	}
	// Valida se o horário solicitado está dentro dos horários disponíveis do dia
	if !agendaDoDia.PermiteAgendamento(input.DataHoraInicio, dataHorarioFim) {
		return nil, ErrHorarioIndisponivel
	}

	// Valida conflitos com outros agendamentos
	// Garante que o prestador não tenha dois atendimentos no mesmo horário
	// Essa regra fica no Aggregate Root (Prestador)
	if !prestador.PodeAgendar(input.DataHoraInicio, dataHorarioFim) {
		return nil, ErrHorarioIndisponivel
	}
	// Um prestador não pode ter dois atendimentos no mesmo período
	conflitosPrestador, err := s.agendamentoRepo.BuscarPorPrestadorEPeriodo(prestador.ID, input.DataHoraInicio, dataHorarioFim)
	if err != nil {
		return nil, err
	}
	if len(conflitosPrestador) > 0 {
		return nil, ErrPrestadorOcupado
	}

	// Um cliente não pode ter dois agendamentos simultâneos
	conflitosCliente, err := s.agendamentoRepo.BuscarPorClienteEPeriodo(cliente.ID, input.DataHoraInicio, dataHorarioFim)
	if err != nil {
		return nil, err
	}
	if len(conflitosCliente) > 0 {
		return nil, ErrClienteOcupado
	}
	agendamento, err := domain.NovoAgendamento(
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
		return nil, err
	}

	out := mapper.NovoAgendamentoOutput(agendamento)
	return out, nil
}
