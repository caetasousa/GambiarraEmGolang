package response

import "meu-servico-agenda/internal/core/application/output"

// LoginResponse representa a resposta de login unificado
type LoginResponse struct {
	Token string      `json:"token"`
	Role  string      `json:"role"`
	User  interface{} `json:"user"`
}

// FromLoginOutput converte output.LoginOutput para LoginResponse
func FromLoginOutput(out *output.LoginOutput) LoginResponse {
	return LoginResponse{
		Token: out.Token,
		Role:  out.Role,
		User:  out.User,
	}
}
