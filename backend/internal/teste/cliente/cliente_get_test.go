package teste

import (
	"encoding/json"
	"net/http"
	"testing"

	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/core/application/output"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCliente_Sucesso(t *testing.T) {
	router, _ := SetupRouterCliente()

	// 1. Cadastrar um cliente
	input := request.ClienteRequest{
		Nome:     "Beatriz",
		Email:    "beatriz@gmail.com",
		Telefone: "62999977848",
		Senha:    "senha123",
	}

	rrPost := SetupPostClienteRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code, "Cliente deveria ser criado com sucesso")

	var clienteCriado output.BuscarClienteOutput
	err := json.Unmarshal(rrPost.Body.Bytes(), &clienteCriado)
	require.NoError(t, err, "Resposta do POST deve ser um JSON válido")

	// 2. Buscar o cliente criado
	rrGet := SetupGetClienteRequest(router, clienteCriado.ID)

	assert.Equal(t, http.StatusOK, rrGet.Code, "Esperado 200 OK para cliente existente")

	var clienteBuscado output.BuscarClienteOutput
	err = json.Unmarshal(rrGet.Body.Bytes(), &clienteBuscado)
	require.NoError(t, err, "Resposta do GET deve ser um JSON válido")

	// 3. Validar dados
	assert.Equal(t, clienteCriado.ID, clienteBuscado.ID, "IDs devem ser iguais")
	assert.Equal(t, clienteCriado.Nome, clienteBuscado.Nome, "Nomes devem ser iguais")
	assert.Equal(t, clienteCriado.Email, clienteBuscado.Email, "Emails devem ser iguais")
	assert.Equal(t, clienteCriado.Telefone, clienteBuscado.Telefone, "Telefones devem ser iguais")
}

func TestGetCliente_NaoEncontrado(t *testing.T) {
	router, _ := SetupRouterCliente()

	rr := SetupGetClienteRequest(router, "id-inexistente")

	assert.Equal(t, http.StatusNotFound, rr.Code, "Esperado 404 para cliente inexistente")

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")

	_, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
}

func TestGetCliente_IdVazio(t *testing.T) {
	router, _ := SetupRouterCliente()

	rr := SetupGetClienteRequest(router, "")

	// Quando o ID está vazio, a rota não é encontrada
	assert.Equal(t, http.StatusNotFound, rr.Code, "Esperado 404 para ID vazio")
}
