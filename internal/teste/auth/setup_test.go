package auth_test

import (
	"meu-servico-agenda/internal/adapters/http/middleware"
	"meu-servico-agenda/internal/core/domain"
	"meu-servico-agenda/internal/infra/auth"

	"github.com/gin-gonic/gin"
)

const (
	TestJWTSecret = "test-secret-key-for-auth-rbac-tests"
)

// SetupAuthRouter cria um roteador Gin com rotas protegidas para teste de RBAC
func SetupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	api := router.Group("/api/v1")
	
	// Rota que requer apenas autenticação
	api.GET("/test/auth", middleware.AuthMiddleware(TestJWTSecret), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "autenticado"})
	})

	// Rota que requer Role Admin
	api.GET("/test/admin", middleware.AuthMiddleware(TestJWTSecret), middleware.RequireRole(domain.RoleAdmin), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "admin ok"})
	})

	// Rota que requer Role Cliente ou Admin
	api.GET("/test/cliente-admin", middleware.AuthMiddleware(TestJWTSecret), middleware.RequireRole(domain.RoleCliente, domain.RoleAdmin), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "cliente/admin ok"})
	})

	// Rota que requer ser Dono (resource id na rota) ou Admin
	api.GET("/test/owner/:id", middleware.AuthMiddleware(TestJWTSecret), middleware.RequireOwnerOrAdmin(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "owner ok"})
	})

	// Rota que requer ser Dono ou Role específica
	api.GET("/test/owner-or-prestador/:id", middleware.AuthMiddleware(TestJWTSecret), middleware.RequireOwnerOrRole(domain.RolePrestador), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "owner or prestador ok"})
	})

	return router
}

// GenerateTestToken gera um token JWT para testes
func GenerateTestToken(userID, email string, role domain.Role) string {
	token, _ := auth.GenerateToken(userID, email, role, TestJWTSecret, 1)
	return token
}

// GenerateExpiredToken gera um token já expirado para testes
func GenerateExpiredToken(userID, email string, role domain.Role) string {
	// Usamos GenerateToken com expiração negativa
	token, _ := auth.GenerateToken(userID, email, role, TestJWTSecret, -1)
	return token
}
