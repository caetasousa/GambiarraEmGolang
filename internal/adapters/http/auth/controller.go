package auth

import (
	"net/http"

	"meu-servico-agenda/internal/adapters/http/auth/request"
	"meu-servico-agenda/internal/adapters/http/auth/response"
	"meu-servico-agenda/internal/core/application/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NovoAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login é o handler para a rota POST /login
// @Summary Faz login de um usuário (cliente, prestador ou admin)
// @Description Recebe email e senha e retorna um token JWT se as credenciais forem válidas
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credentials body request.LoginRequest true "Credenciais de login"
// @Success 200 {object} response.ClienteLoginResponse "Login bem-sucedido, retorna token e dados do usuário"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 401 {object} domain.ErrorResponse "Credenciais inválidas"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Converter para input
	inputData := req.ToLoginInput()

	// Fazer login usando AuthService
	loginOutput, err := ctrl.authService.LoginCliente(inputData)
	if err != nil {
		errorMessage := err.Error()

		// TRATAMENTO DE CREDENCIAIS INVÁLIDAS
		if errorMessage == "credenciais inválidas" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha incorretos"})
			return
		}

		// TRATAMENTO DE ERRO GENÉRICO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer login"})
		return
	}

	// Converter output para response
	resp := response.FromClienteLoginOutput(loginOutput)
	c.JSON(http.StatusOK, resp)
}

// LoginPrestador é o handler para login de prestador
// @Summary Faz login de um prestador
// @Description Recebe email e senha e retorna um token JWT se as credenciais forem válidas
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credentials body request.LoginRequest true "Credenciais de login"
// @Success 200 {object} response.PrestadorLoginResponse "Login bem-sucedido, retorna token e dados do prestador"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos"
// @Failure 401 {object} domain.ErrorResponse "Credenciais inválidas"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /prestadores/login [post]
func (ctrl *AuthController) LoginPrestador(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Converter para input
	inputData := req.ToLoginInput()

	// Fazer login usando AuthService
	loginOutput, err := ctrl.authService.LoginPrestador(inputData)
	if err != nil {
		errorMessage := err.Error()

		// TRATAMENTO DE CREDENCIAIS INVÁLIDAS
		if errorMessage == "credenciais inválidas" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou senha incorretos"})
			return
		}

		// TRATAMENTO DE ERRO GENÉRICO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao fazer login"})
		return
	}

	// Converter output para response
	resp := response.FromPrestadorLoginOutput(loginOutput)
	c.JSON(http.StatusOK, resp)
}
