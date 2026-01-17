package mapper

import (
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/domain"
)

// ClienteToLoginData converte um Cliente de domínio para LoginClienteData (SEM senha)
func ClienteToLoginData(cliente *domain.Cliente) output.LoginClienteData {
	return output.LoginClienteData{
		ID:       cliente.ID,
		Nome:     cliente.Nome,
		Email:    cliente.Email,
		Telefone: cliente.Telefone,
	}
}

// PrestadorToLoginData converte um Prestador de domínio para LoginPrestadorData (SEM senha)
func PrestadorToLoginData(prestador *domain.Prestador) output.LoginPrestadorData {
	return output.LoginPrestadorData{
		ID:        prestador.ID,
		Nome:      prestador.Nome,
		Email:     prestador.Email,
		Telefone:  prestador.Telefone,
		Cpf:       prestador.Cpf,
		ImagemUrl: prestador.ImagemUrl,
	}
}

// ClienteToLoginOutput converte um Cliente de domínio para LoginClienteOutput completo (com token)
func ClienteToLoginOutput(cliente *domain.Cliente, token string) *output.LoginClienteOutput {
	return &output.LoginClienteOutput{
		Token: token,
		User:  ClienteToLoginData(cliente),
	}
}

// PrestadorToLoginOutput converte um Prestador de domínio para LoginPrestadorOutput completo (com token)
func PrestadorToLoginOutput(prestador *domain.Prestador, token string) *output.LoginPrestadorOutput {
	return &output.LoginPrestadorOutput{
		Token: token,
		User:  PrestadorToLoginData(prestador),
	}
}
