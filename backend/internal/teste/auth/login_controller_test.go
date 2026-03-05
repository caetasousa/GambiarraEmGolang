package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"meu-servico-agenda/internal/adapters/http/auth/request"
	"meu-servico-agenda/internal/adapters/http/auth/response"
	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginController(t *testing.T) {
	router, clienteRepo, prestadorRepo := SetupLoginController()
	SeedData(clienteRepo, prestadorRepo)

	t.Run("Deve fazer login como Cliente com sucesso", func(t *testing.T) {
		loginReq := request.LoginRequest{
			Email: "cliente@test.com",
			Senha: "senha123",
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp response.LoginResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, domain.RoleCliente.String(), resp.Role)
		
		// Verificar se os dados do usuário estão presentes
		userData := resp.User.(map[string]interface{})
		assert.Equal(t, "cliente@test.com", userData["email"])
		assert.Equal(t, "Cliente Teste", userData["nome"])
	})

	t.Run("Deve fazer login como Prestador com sucesso", func(t *testing.T) {
		loginReq := request.LoginRequest{
			Email: "prestador@test.com",
			Senha: "senha123",
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp response.LoginResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, domain.RolePrestador.String(), resp.Role)
		
		// Verificar se os dados do prestador estão presentes
		userData := resp.User.(map[string]interface{})
		assert.Equal(t, "prestador@test.com", userData["email"])
		assert.Equal(t, "Prestador Teste", userData["nome"])
		assert.NotEmpty(t, userData["cpf"])
	})

	t.Run("Deve fazer login como Admin com sucesso", func(t *testing.T) {
		loginReq := request.LoginRequest{
			Email: "admin@test.com",
			Senha: "senha123",
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp response.LoginResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, domain.RoleAdmin.String(), resp.Role)
	})

	t.Run("Deve retornar 401 para senha incorreta", func(t *testing.T) {
		loginReq := request.LoginRequest{
			Email: "cliente@test.com",
			Senha: "senha_errada",
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Deve retornar 401 para usuário inexistente", func(t *testing.T) {
		loginReq := request.LoginRequest{
			Email: "nao_existe@test.com",
			Senha: "senha123",
		}
		body, _ := json.Marshal(loginReq)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
