package service

import (
	"meu-servico-agenda/internal/adapters/http/auth/request"
	"meu-servico-agenda/internal/core/application/input"
	"meu-servico-agenda/internal/core/application/mapper"
	"meu-servico-agenda/internal/core/application/output"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/domain"
	"meu-servico-agenda/internal/infra/auth"
)

// AuthService gerencia a autenticação e geração de tokens JWT
type AuthService struct {
	jwtSecret        string
	jwtExpirationHrs int
	clienteRepo      port.ClienteRepositorio
	prestadorRepo    port.PrestadorRepositorio
}

// NovoAuthService cria uma nova instância do serviço de autenticação
func NovoAuthService(jwtSecret string, jwtExpirationHrs int, clienteRepo port.ClienteRepositorio, prestadorRepo port.PrestadorRepositorio) *AuthService {
	return &AuthService{
		jwtSecret:        jwtSecret,
		jwtExpirationHrs: jwtExpirationHrs,
		clienteRepo:      clienteRepo,
		prestadorRepo:    prestadorRepo,
	}
}

// GenerateToken gera um token JWT para o usuário autenticado
func (s *AuthService) GenerateToken(userID, email string, role domain.Role) (string, error) {
	token, err := auth.GenerateToken(userID, email, role, s.jwtSecret, s.jwtExpirationHrs)
	if err != nil {
		return "", ErrFalhaAoGerarToken
	}
	return token, nil
}

// ValidateToken valida um token JWT e retorna os claims
func (s *AuthService) ValidateToken(tokenString string) (*auth.JWTClaims, error) {
	claims, err := auth.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// LoginCliente autentica especificamente um cliente
func (s *AuthService) LoginCliente(inputData *input.LoginInput) (*output.LoginClienteOutput, error) {
	// Buscar cliente por email
	cliente, err := s.clienteRepo.BuscarPorEmail(inputData.Email)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	if cliente == nil {
		return nil, ErrCredenciaisInvalidas
	}

	// Validar senha
	if err := request.ValidarSenha(cliente.PasswordHash, inputData.Senha); err != nil {
		return nil, ErrCredenciaisInvalidas
	}

	// Gerar JWT token
	token, err := s.GenerateToken(cliente.ID, cliente.Email, cliente.Role)
	if err != nil {
		return nil, ErrFalhaAoGerarToken
	}

	// Retornar token e dados do cliente
	return mapper.ClienteToLoginOutput(cliente, token), nil
}

// LoginPrestador autentica especificamente um prestador
func (s *AuthService) LoginPrestador(inputData *input.LoginInput) (*output.LoginPrestadorOutput, error) {
	// Buscar prestador por email
	prestador, err := s.prestadorRepo.BuscarPorEmail(inputData.Email)
	if err != nil {
		return nil, ErrFalhaInfraestrutura
	}

	if prestador == nil {
		return nil, ErrCredenciaisInvalidas
	}

	// Validar senha
	if err := request.ValidarSenha(prestador.PasswordHash, inputData.Senha); err != nil {
		return nil, ErrCredenciaisInvalidas
	}

	// Gerar JWT token
	token, err := s.GenerateToken(prestador.ID, prestador.Email, prestador.Role)
	if err != nil {
		return nil, ErrFalhaAoGerarToken
	}

	// Retornar token e dados do prestador
	return mapper.PrestadorToLoginOutput(prestador, token), nil
}
