package auth_test

import (
	"meu-servico-agenda/internal/adapters/http/auth"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SetupLoginController configura o ambiente para testes de login
func SetupLoginController() (*gin.Engine, port.ClienteRepositorio, port.PrestadorRepositorio) {
	gin.SetMode(gin.TestMode)

	clienteRepo := repository.NewFakeClienteRepositorio()
	catalogoRepo := repository.NovoCatalogoFakeRepo()
	prestadorRepo := repository.NovoFakePrestadorRepositorio(catalogoRepo)

	authService := service.NovoAuthService(TestJWTSecret, 24, clienteRepo, prestadorRepo)
	authController := auth.NovoAuthController(authService)

	router := gin.Default()
	router.POST("/login", authController.Login)

	return router, clienteRepo, prestadorRepo
}

// HashPassword helper para gerar hash de senha para os testes
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// SeedData insere dados iniciais para os testes
func SeedData(clienteRepo port.ClienteRepositorio, prestadorRepo port.PrestadorRepositorio) {
	passwordHash := HashPassword("senha123")

	// 1. Cliente
	cliente := domain.NovoCliente("cliente-1", "Cliente Teste", "cliente@test.com", "11999999999", passwordHash)
	clienteRepo.Salvar(cliente)

	// 2. Prestador
	prestador, _ := domain.NovoPrestador("prestador-1", "Prestador Teste", "04423258196", "prestador@test.com", "11888888888", passwordHash, "", nil)
	prestadorRepo.Salvar(prestador)

	// 3. Admin (usando o role de Admin no domínio de cliente ou conforme definido no sistema)
	admin := domain.NovoCliente("admin-1", "Admin Teste", "admin@test.com", "11777777777", passwordHash)
	admin.Role = domain.RoleAdmin
	clienteRepo.Salvar(admin)
}
