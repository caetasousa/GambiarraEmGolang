package response

import "meu-servico-agenda/internal/core/application/output"

// ClienteLoginResponse representa a resposta de login unificado (cliente ou prestador)
type ClienteLoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"` // Pode ser LoginClienteData ou LoginPrestadorData
}

// FromClienteLoginOutput converte output.LoginOutput para ClienteLoginResponse
func FromClienteLoginOutput(out *output.LoginClienteOutput) ClienteLoginResponse {
	return ClienteLoginResponse{
		Token: out.Token,
		User:  out.User,
	}
}

// PrestadorLoginResponse representa a resposta de login de prestador
type PrestadorLoginResponse struct {
	Token string                     `json:"token"`
	User  output.LoginPrestadorData `json:"user"`
}

// FromPrestadorLoginOutput converte output.LoginPrestadorOutput para PrestadorLoginResponse
func FromPrestadorLoginOutput(out *output.LoginPrestadorOutput) PrestadorLoginResponse {
	return PrestadorLoginResponse{
		Token: out.Token,
		User:  out.User,
	}
}
