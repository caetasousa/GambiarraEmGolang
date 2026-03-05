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

func TestPostCliente_Sucesso(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Ana Silva",
		Email:    "ana@example.com",
		Telefone: "62999697481",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code, "Não deveria retornar 400, JSON é válido")
	assert.NotEqual(t, http.StatusInternalServerError, rr.Code, "Não deveria causar erro interno")
	assert.Equal(t, http.StatusCreated, rr.Code, "Esperado 201 Created")

	var resp output.BuscarClienteOutput
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err, "Resposta deve ser um JSON válido")

	assert.Equal(t, input.Nome, resp.Nome)
	assert.Equal(t, input.Email, resp.Email)
	assert.Equal(t, input.Telefone, resp.Telefone)
	assert.NotZero(t, resp.ID, "O serviço deve gerar ID")
}

func TestPostCliente_EmailInvalido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Carlos",
		Email:    "email-invalido",
		Telefone: "62999697482",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para email inválido")
	
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")
	
	_, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Contains(t, resp["error"], "Email", "Mensagem de erro deve mencionar Email")
}

func TestPostCliente_NomeRequerido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "",
		Email:    "email@valido.com",
		Telefone: "62999697482",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para nome vazio")
	
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")
	
	_, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Contains(t, resp["error"], "Nome", "Mensagem de erro deve mencionar Nome")
}

func TestPostCliente_TelefoneRequerido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Eduardo",
		Email:    "email@valido.com",
		Telefone: "",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para telefone vazio")
	
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")
	
	_, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Contains(t, resp["error"], "Telefone", "Mensagem de erro deve mencionar Telefone")
}

func TestPostCliente_SenhaRequerida(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Fernando",
		Email:    "fernando@email.com",
		Telefone: "62999697483",
		Senha:    "",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para senha vazia")
	
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")
	
	_, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Contains(t, resp["error"], "Senha", "Mensagem de erro deve mencionar Senha")
}

func TestPostCliente_NomeMuitoCurto(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "AB",
		Email:    "ab@example.com",
		Telefone: "62999974848",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para nome com < 3 caracteres")
}

func TestPostCliente_TelefoneMuitoCurto(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "João",
		Email:    "joao@example.com",
		Telefone: "1234567",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para telefone com < 8 dígitos")
}

func TestPostCliente_SenhaMuitoCurta(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Maria",
		Email:    "maria@example.com",
		Telefone: "62999974848",
		Senha:    "123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para senha com < 8 caracteres")
}

func TestPostCliente_TelefoneComTipoInvalido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := map[string]interface{}{
		"nome":     "Carlos",
		"telefone": 6299974848, // número em vez de string
		"email":    "carlos@example.com",
		"senha":    "senha123",
	}

	rr := PostJSONFromMap(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para telefone com tipo inválido")
}

func TestPostCliente_EmailOmitido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Maria",
		Email:    "",
		Telefone: "62999974848",
		Senha:    "senha123",
	}

	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para email omitido")
}

func TestPostCliente_EmailDuplicado(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Primeiro Cliente",
		Email:    "duplicado@email.com",
		Telefone: "62999974848",
		Senha:    "senha123",
	}

	// Primeiro cadastro - deve funcionar
	rr1 := SetupPostClienteRequest(router, input)
	require.Equal(t, http.StatusCreated, rr1.Code, "Primeiro cadastro deve funcionar")

	// Segundo cadastro com mesmo email - deve falhar
	input2 := request.ClienteRequest{
		Nome:     "Segundo Cliente",
		Email:    "duplicado@email.com",
		Telefone: "62999974849",
		Senha:    "senha456",
	}

	rr2 := SetupPostClienteRequest(router, input2)
	assert.Equal(t, http.StatusConflict, rr2.Code, "Esperado 409 Conflict para email duplicado")
}
