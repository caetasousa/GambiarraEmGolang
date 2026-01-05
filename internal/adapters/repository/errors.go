package repository

import "errors"


var (
	// Erros genéricos de persistência
	ErrNaoEncontrado      = errors.New("registro não encontrado")
	ErrFalhaAoSalvar      = errors.New("falha ao salvar registro")
	ErrFalhaAoAtualizar   = errors.New("falha ao atualizar registro")
	ErrFalhaAoDeletar     = errors.New("falha ao deletar registro")
	ErrFalhaAoListar      = errors.New("falha ao listar registros")
	ErrFalhaDeTransacao   = errors.New("falha na transação")

	// Erros de constraint/duplicidade
	ErrEmailDuplicado     = errors.New("email já cadastrado")
	ErrCPFDuplicado       = errors.New("cpf já cadastrado")
	ErrRegistroDuplicado  = errors.New("registro duplicado")

	// Erros específicos por entidade (opcional, se quiser mais granularidade)
	ErrCatalogoNaoEncontrado   = errors.New("catálogo não encontrado")
	ErrPrestadorNaoEncontrado  = errors.New("prestador não encontrado")
	ErrClienteNaoEncontrado    = errors.New("cliente não encontrado")
	ErrAgendaNaoEncontrada     = errors.New("agenda não encontrada")
	ErrAgendamentoNaoEncontrado = errors.New("agendamento não encontrado")
)