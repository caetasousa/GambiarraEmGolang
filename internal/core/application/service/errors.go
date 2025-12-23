package service

import "errors"

var (
	//validação de cpf
	ErrCPFJaCadastrado = errors.New("cpf já possui um cadastro")

	//validação de catálogo
	ErrCatalogoNaoExiste     = errors.New("catálogo não existe")
	ErrCatalogoNaoEncontrado = errors.New("catálogo não encontrado")
	ErrCatalogoInvalido      = errors.New("catálogo inválido")

	ErrFalhaInfraestrutura = errors.New("falha na infraestrutura")

	//validação de prestador
	ErrPrestadorNaoEncontrado = errors.New("prestador não encontrado")
	ErrPrestadorInvalido      = errors.New("prestador inválido")
	ErrPrestadorNaoExiste     = errors.New("prestador não encontrado")
	ErrPrestadorOcupado       = errors.New("prestador já possui agendamento neste horário")
	ErrPrestadorInativo       = errors.New("prestador está inativo")

	//validação de cliente
	ErrClienteNaoEncontrado = errors.New("cliente não encontrado")
	ErrAoSalvarCliente      = errors.New("falha ao salvar cliente: ")
	ErrClienteInvalido      = errors.New("cliente inválido")
	ErrClienteNaoExiste     = errors.New("cliente não encontrado")
	ErrClienteOcupado       = errors.New("cliente já possui agendamento neste horário")

	//validação de agendamento
	ErrDataHoraInvalida    = errors.New("data/hora de agendamento inválida")
	ErrHorarioIndisponivel = errors.New("horário indisponível")
	ErrDiaIndisponivel     = errors.New("dia indisponível para agendamentos")
	ErrAgendaDuplicada     = errors.New("Agenda diaria duplicada")
)
