package service

import (
	"errors"
	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
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

var (
	ErrClienteInvalido   = errors.New("cliente inválido")
	ErrPrestadorInvalido = errors.New("prestador inválido")
	ErrCatalogoInvalido  = errors.New("catálogo inválido")

	ErrClienteNaoExiste   = errors.New("cliente não encontrado")
	ErrPrestadorNaoExiste = errors.New("prestador não encontrado")

	ErrDataHoraInvalida    = errors.New("data/hora de agendamento inválida")
	ErrHorarioIndisponivel = errors.New("horário indisponível")
	ErrDiaIndisponivel     = errors.New("dia indisponível para agendamentos")

	ErrPrestadorOcupado = errors.New("prestador já possui agendamento neste horário")
	ErrClienteOcupado   = errors.New("cliente já possui agendamento neste horário")
)

func (s *AgendamentoService) CadastraAgendamento(input request_agendamento.AgendamentoRequest) (*domain.Agendamento, error) {
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

	dataHoraInicio, err := time.Parse(time.RFC3339, input.DataHoraInicio)
	if err != nil {
		return nil, ErrDataHoraInvalida
	}

	dataHorarioFim := dataHoraInicio.Add(time.Duration(catalogo.DuracaoPadrao) * time.Minute)

	// Busca a agenda do prestador para o dia solicitado
	// Aqui validamos se o prestador trabalha nesse dia específico
	agendaDoDia, err := s.prestadorRepo.BuscarAgendaDoDia(prestador.ID, dataHoraInicio.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	// Valida se o dia é atendido pelo prestador
	// Exemplo: prestador não trabalha à tarde ou não trabalha nesse dia
	if agendaDoDia == nil {
		return nil, ErrDiaIndisponivel
	}
	// Valida se o horário solicitado está dentro dos horários disponíveis do dia
	if !agendaDoDia.PermiteAgendamento(dataHoraInicio, dataHorarioFim) {
		return nil, ErrHorarioIndisponivel
	}

	// Valida conflitos com outros agendamentos
	// Garante que o prestador não tenha dois atendimentos no mesmo horário
	// Essa regra fica no Aggregate Root (Prestador)
	if !prestador.PodeAgendar(dataHoraInicio, dataHorarioFim) {
		return nil, ErrHorarioIndisponivel
	}
	// 6️⃣ Conflito: PRESTADOR
	conflitosPrestador, err := s.agendamentoRepo.BuscarPorPrestadorEPeriodo(prestador.ID, dataHoraInicio, dataHorarioFim)
	if err != nil {
		return nil, err
	}
	if len(conflitosPrestador) > 0 {
		return nil, ErrPrestadorOcupado
	}

	// 7️⃣ Conflito: CLIENTE
	conflitosCliente, err := s.agendamentoRepo.BuscarPorClienteEPeriodo(cliente.ID, dataHoraInicio, dataHorarioFim)
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
		dataHoraInicio,
		dataHorarioFim,
		input.Notas,
	)
	if err != nil {
		return nil, err
	}

	return agendamento, nil
}
