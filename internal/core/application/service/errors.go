package service

import "errors"

var (
	// Validação de CPF
	ErrCPFJaCadastrado = errors.New("cpf já possui um cadastro")

	// Validação de catálogo
	ErrCatalogoNaoEncontrado = errors.New("catálogo não encontrado")
	ErrCatalogoInvalido      = errors.New("catálogo inválido")

	// Erro genérico de infraestrutura
	ErrFalhaInfraestrutura = errors.New("falha na infraestrutura")

	// Validação de prestador
	ErrPrestadorNaoEncontrado         = errors.New("prestador não encontrado")
	ErrPrestadorInvalido              = errors.New("prestador inválido")
	ErrPrestadorOcupado               = errors.New("prestador já possui agendamento neste horário")
	ErrPrestadorInativo               = errors.New("prestador está inativo")
	ErrFormatoDataInvalido            = errors.New("formato de data inválido. Use YYYY-MM-DD")
	ErrAoBuscarPrestadoresDisponiveis = errors.New("erro ao buscar prestadores disponíveis")
	ErrAoContarPrestadoresDisponiveis = errors.New("erro ao contar prestadores disponíveis")

	// Validação de cliente
	ErrClienteNaoEncontrado = errors.New("cliente não encontrado")
	ErrClienteInvalido      = errors.New("cliente inválido")
	ErrClienteOcupado       = errors.New("cliente já possui agendamento neste horário")

	// Validação de agendamento
	ErrDataHoraInvalida    = errors.New("data/hora de agendamento inválida")
	ErrHorarioIndisponivel = errors.New("horário indisponível")
	ErrDiaIndisponivel     = errors.New("dia indisponível para agendamentos")
	ErrAgendaNaoEncontrada = errors.New("agenda não encontrada")
	ErrAgendamentoDuplo    = errors.New("já existe um agendamento para essa categoria neste dia")
)
