package http

import "errors"

var (
	// Erros de data/hora
	ErrFormatoDataHoraInvalido = errors.New("formato de data/hora inválido")
	ErrDataInvalida            = errors.New("data inválida")
	ErrHoraInicioInvalida      = errors.New("hora_inicio inválida")
	ErrHoraFimInvalida         = errors.New("hora_fim inválida")

	// Erros de validação
	ErrCPFInvalido = errors.New("cpf inválido")
)