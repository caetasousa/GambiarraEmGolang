package teste

import (
	"encoding/json"
	"net/http"
	"testing"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrestador_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria Silva",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		ImagemUrl:   "https://exemplo.com/img-prestador.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err := json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)
	require.NotEmpty(t, createResp.ID)

	// Buscar prestador criado
	rrGet := SetupGetPrestadorRequest(router, createResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	var prestadorResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrGet.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	// Validações
	assert.Equal(t, createResp.ID, prestadorResp.ID)
	assert.Equal(t, "Maria Silva", prestadorResp.Nome)
	assert.Equal(t, "04423258196", prestadorResp.Cpf)
	assert.Equal(t, "maria@email.com", prestadorResp.Email)
	assert.True(t, prestadorResp.Ativo)
	assert.Len(t, prestadorResp.Catalogo, 1)
	assert.Empty(t, prestadorResp.Agenda)
}

func TestGetPrestador_NaoEncontrado(t *testing.T) {
	router, _ := SetupPostPrestador()
	rr := SetupGetPrestadorRequest(router, "id-inexistente")
	require.Equal(t, http.StatusNotFound, rr.Code)
}