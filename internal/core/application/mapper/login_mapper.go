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
		Cpf:       prestador.Cpf.String(),
		ImagemUrl: prestador.ImagemUrl,
	}
}

// ClienteToLoginOutput converte um Cliente de domínio para LoginOutput completo (com token)
func ClienteToLoginOutput(cliente *domain.Cliente, token string) *output.LoginOutput {
	return &output.LoginOutput{
		Token: token,
		Role:  cliente.Role.String(),
		User:  ClienteToLoginData(cliente),
	}
}

// PrestadorToLoginOutput converte um Prestador de domínio para LoginOutput completo (com token)
func PrestadorToLoginOutput(prestador *domain.Prestador, token string) *output.LoginOutput {
	return &output.LoginOutput{
		Token: token,
		Role:  prestador.Role.String(),
		User:  PrestadorToLoginData(prestador),
	}
}
