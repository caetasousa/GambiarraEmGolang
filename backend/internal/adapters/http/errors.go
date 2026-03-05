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

	// Erros de Autenticação e Autorização
	ErrTokenNaoFornecido     = errors.New("token de autenticação não fornecido")
	ErrTokenInvalido         = errors.New("formato de token inválido. Use: Bearer <token>")
	ErrUsuarioNaoAutenticado = errors.New("usuário não autenticado")
	ErrPermissaoInsuficiente = errors.New("acesso negado: permissão insuficiente")
	ErrAcessoSomenteProprio  = errors.New("acesso negado: você só pode acessar seus próprios dados")
)