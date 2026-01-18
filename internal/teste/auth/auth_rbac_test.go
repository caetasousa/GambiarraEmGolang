package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	adapterhttp "meu-servico-agenda/internal/adapters/http"
	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	router := SetupAuthRouter()

	t.Run("Deve falhar quando token não é fornecido", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/auth", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), adapterhttp.ErrTokenNaoFornecido.Error())
	})

	t.Run("Deve falhar quando formato do token é inválido", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/auth", nil)
		req.Header.Set("Authorization", "Invalido token123")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "formato de token inválido")
	})

	t.Run("Deve falhar quando token está expirado", func(t *testing.T) {
		token := GenerateExpiredToken("user-123", "test@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/auth", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		// A mensagem de erro vem do ValidateToken -> auth.ErrTokenExpirado
		// De acordo com middleware/auth.go:39: c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		assert.Contains(t, rr.Body.String(), "token expirado")
	})

	t.Run("Deve permitir quando token é válido", func(t *testing.T) {
		token := GenerateTestToken("user-123", "test@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/auth", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "autenticado")
	})
}

func TestRBAC_RequireRole(t *testing.T) {
	router := SetupAuthRouter()

	t.Run("Admin deve acessar rota de Admin", func(t *testing.T) {
		token := GenerateTestToken("admin-1", "admin@test.com", domain.RoleAdmin)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Cliente não deve acessar rota de Admin", func(t *testing.T) {
		token := GenerateTestToken("cliente-1", "cliente@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), adapterhttp.ErrPermissaoInsuficiente.Error())
	})

	t.Run("Cliente e Admin devem acessar rota compartilhada", func(t *testing.T) {
		// Cliente
		tokenCli := GenerateTestToken("cliente-1", "cliente@test.com", domain.RoleCliente)
		req1, _ := http.NewRequest(http.MethodGet, "/api/v1/test/cliente-admin", nil)
		req1.Header.Set("Authorization", "Bearer "+tokenCli)
		rr1 := httptest.NewRecorder()
		router.ServeHTTP(rr1, req1)
		assert.Equal(t, http.StatusOK, rr1.Code)

		// Admin
		tokenAdm := GenerateTestToken("admin-1", "admin@test.com", domain.RoleAdmin)
		req2, _ := http.NewRequest(http.MethodGet, "/api/v1/test/cliente-admin", nil)
		req2.Header.Set("Authorization", "Bearer "+tokenAdm)
		rr2 := httptest.NewRecorder()
		router.ServeHTTP(rr2, req2)
		assert.Equal(t, http.StatusOK, rr2.Code)
	})
}

func TestRBAC_RequireOwnerOrAdmin(t *testing.T) {
	router := SetupAuthRouter()

	t.Run("Usuário acessando seu próprio recurso", func(t *testing.T) {
		token := GenerateTestToken("user-123", "user@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/owner/user-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Admin acessando recurso de qualquer usuário", func(t *testing.T) {
		token := GenerateTestToken("admin-1", "admin@test.com", domain.RoleAdmin)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/owner/user-outro", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Usuário acessando recurso de outro usuário deve falhar", func(t *testing.T) {
		token := GenerateTestToken("user-123", "user@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/owner/user-outro", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), adapterhttp.ErrAcessoSomenteProprio.Error())
	})
}

func TestRBAC_RequireOwnerOrRole(t *testing.T) {
	router := SetupAuthRouter()

	t.Run("Usuário acessando seu recurso (sem ter a role específica)", func(t *testing.T) {
		token := GenerateTestToken("user-123", "user@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/owner-or-prestador/user-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Usuário com a role específica acessando recurso de outro", func(t *testing.T) {
		token := GenerateTestToken("prestador-1", "prestador@test.com", domain.RolePrestador)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/owner-or-prestador/user-123", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Usuário sem a role e não sendo dono deve falhar", func(t *testing.T) {
		token := GenerateTestToken("cliente-1", "cliente@test.com", domain.RoleCliente)
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/test/owner-or-prestador/prestador-1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), adapterhttp.ErrPermissaoInsuficiente.Error())
	})
}
