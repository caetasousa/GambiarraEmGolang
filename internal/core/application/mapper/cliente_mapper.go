package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

// ClienteToOutput converte um Cliente de domínio para output (SEM senha)
func ClienteToOutput(cliente *domain.Cliente) output.BuscarClienteOutput {
	return output.BuscarClienteOutput{
		ID:       cliente.ID,
		Nome:     cliente.Nome,
		Email:    cliente.Email,
		Telefone: cliente.Telefone,
		Role:     string(cliente.Role),
	}
}
